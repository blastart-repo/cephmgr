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
		Run: func(cmd *cobra.Command, _ []string) {

			user := &User{
				ID:          userID,
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
	createCmd.PersistentFlags().StringVarP(&userID, "user", "u", "", "Ceph user ID")

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
