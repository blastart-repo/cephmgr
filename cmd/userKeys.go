package cmd

import (
	"context"
	"fmt"

	"github.com/ceph/go-ceph/rgw/admin"
	"github.com/spf13/cobra"
)

// keysCmd represents the keys command
var (
	keysCmd = &cobra.Command{
		Use:   "keys",
		Short: "User keys operations",
		Long:  `User keys operations`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
	removeKeyCmd = &cobra.Command{
		Use:   "remove",
		Short: "Remove user keys",
		Long:  `Remove user keys`,
		Args:  cobra.ExactArgs(2), // Require exactly 2 arguments (UID and AccessKey)
		Run: func(cmd *cobra.Command, args []string) {
			uid := args[0]       // Get the UID from the command line argument
			accessKey := args[1] // Get the AccessKey from the command line argument

			err := removeUserKeys(uid, accessKey)
			if err != nil {
				fmt.Println(err)
				cmd.Help()
			}
		},
	}
	addKeyCmd = &cobra.Command{
		Use:   "add",
		Short: "Add new key to user",
		Long:  `ToDo`,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			uid := args[0] // Get the UID from the command line argument

			err := addUserKey(uid)
			if err != nil {
				fmt.Println(err)
				cmd.Help()
			}
		},
	}
)

func init() {
	userCmd.AddCommand(keysCmd)
	keysCmd.AddCommand(removeKeyCmd)
	keysCmd.AddCommand(addKeyCmd)
}

func removeUserKeys(uid string, accessKey string) error {
	c, err := admin.New(cephHost, cephAccessKey, cephAccessSecret, nil)
	if err != nil {
		return err
	}

	err = c.RemoveKey(context.Background(), admin.UserKeySpec{UID: uid, KeyType: "s3", AccessKey: accessKey})

	if err != nil {
		return err
	}

	fmt.Printf("Successfully removed keys for user %s with AccessKey %s\n", uid, accessKey)
	return nil
}

func addUserKey(uid string) error {
	c, err := admin.New(cephHost, cephAccessKey, cephAccessSecret, nil)
	if err != nil {
		return err
	}

	// Create a pointer to a boolean with the value true
	generateKey := true

	resp, err := c.CreateKey(context.Background(), admin.UserKeySpec{UID: uid, KeyType: "s3", GenerateKey: &generateKey})

	if err != nil {
		return err
	}
	fmt.Println(resp)
	fmt.Printf("Successfully added new keys for user %s\n", uid)
	return nil
}
