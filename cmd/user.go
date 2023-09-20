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
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/ceph/go-ceph/rgw/admin"
	"github.com/spf13/cobra"
)

var (
	userCmd = &cobra.Command{
		Use:   "user",
		Short: "Ceph users operations",
		Long:  `Get users information. Create new users. Change users caps`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
	getuserCmd = &cobra.Command{
		Use:   "get",
		Short: "Get user info",
		Long:  `Get user info`,
		Run: func(cmd *cobra.Command, args []string) {

			user := &User{
				ID:          userName,
				DisplayName: userFullname,
				Email:       userEmail,
			}

			if user.ID == "" {
				cmd.Help()
				os.Exit(1)
			}
			err := getUser(*user)
			if err != nil {
				fmt.Println(err)
				cmd.Help()
			}
		},
	}
	listCmd = &cobra.Command{
		Use:   "list",
		Short: "Get a list of users",
		Long:  `get list of users from the cluster.`,
		Run: func(cmd *cobra.Command, args []string) {
			err := listUsers()
			if err != nil {
				fmt.Println(err)
				cmd.Help()
			}
		},
	}
	deleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "Delete user",
		Long:  `Delete user`,
		Run: func(cmd *cobra.Command, args []string) {

			user := &User{
				ID:          userName,
				DisplayName: userFullname,
				Email:       userEmail,
			}

			if user.ID == "" {
				cmd.Help()
				os.Exit(1)
			}

			err := deleteUser(*user)
			if err != nil {
				fmt.Println(err)
				cmd.Help()
			}
		},
	}
)

func init() {
	rgwCmd.AddCommand(userCmd)

	userCmd.AddCommand(getuserCmd)
	userCmd.AddCommand(listCmd)
	userCmd.AddCommand(deleteCmd)

	userCmd.PersistentFlags().StringVarP(&userName, "user", "u", "", "Ceph user name")
	userCmd.PersistentFlags().StringVarP(&userCaps, "caps", "", "", "User capabilities")
	getuserCmd.MarkFlagRequired("user")
	deleteCmd.MarkFlagRequired("user")
}

func getUser(user User) error {
	c, err := admin.New(cephHost, cephAccessKey, cephAccessSecret, nil)
	if err != nil {
		return err
	}

	u, err := c.GetUser(context.Background(), admin.User{ID: user.ID})

	if err != nil {
		return err
	}

	w := tabwriter.NewWriter(os.Stdout, 10, 1, 5, ' ', 0)

	fs := "%s\t%s\t%s\t%v\n"
	fmt.Fprintln(w, "UID\tFull Name\tEmail\tCaps")
	fmt.Fprintf(w, fs, u.ID, u.DisplayName, u.Email, u.Caps)
	w.Flush()
	return nil
}

func listUsers() error {

	c, err := admin.New(cephHost, cephAccessKey, cephAccessSecret, nil)
	if err != nil {
		return err
	}
	users, err := c.GetUsers(context.Background())

	if err != nil {
		return err
	}

	for _, j := range *users {
		fmt.Println(j)
	}
	return nil
}

func deleteUser(user User) error {

	c, err := admin.New(cephHost, cephAccessKey, cephAccessSecret, nil)
	if err != nil {
		return err
	}

	err = c.RemoveUser(context.Background(), admin.User{ID: user.ID})

	if err != nil {
		return err
	}
	return nil
}
