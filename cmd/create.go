/*
Copyright © 2022 Tarmo Katmuk <tarmo.katmuk@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ceph/go-ceph/rgw/admin"
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var (
	createCmd = &cobra.Command{
		Use:   "create",
		Short: "Create new user",
		Long: `Create new user.
You can also provide capabilities for user with --caps flag:

--caps "buckets=*"`,
		Run: func(cmd *cobra.Command, args []string) {

			user := &User{
				ID:          userName,
				DisplayName: userFullname,
				Email:       userEmail,
				UserCaps:    userCaps,
			}

			err := createUser(*user)
			if err != nil {
				fmt.Println(err)
				cmd.Help()
			}

		},
	}
)

func init() {
	userCmd.AddCommand(createCmd)

	createCmd.Flags().StringVarP(&userFullname, "fullname", "f", "", "Ceph user name")
	createCmd.Flags().StringVarP(&userEmail, "email", "e", "", "Ceph user name")

	createCmd.MarkFlagRequired("user")

}

func createUser(user User) error {

	c, err := admin.New(cephHost, cephAccessKey, cephAccessSecret, nil)
	if err != nil {
		return err
	}
	users, err := c.CreateUser(context.Background(), admin.User{ID: user.ID, DisplayName: user.DisplayName, UserCaps: user.UserCaps})

	if err != nil {
		return err
	}

	buser, _ := json.Marshal(users)

	var userdata User
	_ = json.Unmarshal([]byte(buser), &userdata)

	fmt.Printf("Created user for %s\n", userdata.DisplayName)

	for _, ad := range userdata.Keys {
		fmt.Println("ID:", ad.User)
		fmt.Println("accesskey:", ad.AccessKey)
		fmt.Println("secret:", ad.SecretKey)
	}
	return nil
}
