package cmd

import (
	"context"
	"fmt"

	"github.com/ceph/go-ceph/rgw/admin"
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var (
	deleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "Delete user",
		Long:  `Delete user`,
		Args:  cobra.ExactArgs(1), // Require exactly 1 argument (UID)
		Run: func(cmd *cobra.Command, args []string) {
			user := &User{
				ID: args[0], // Use the first argument as the UID

			}
			err := deleteUser(*user)
			if err != nil {
				fmt.Println(err)
				cmd.Help()
			} else {
				fmt.Printf("Deleted user %s\n", user.ID)
			}

		},
	}
)

func init() {
	userCmd.AddCommand(deleteCmd)

	deleteCmd.MarkFlagRequired("user")

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
