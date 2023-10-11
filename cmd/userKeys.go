package cmd

import (
	"context"
	"fmt"

	"github.com/ceph/go-ceph/rgw/admin"
	"github.com/spf13/cobra"
)

var (
	keysCmd = &cobra.Command{
		Use:   "keys",
		Short: "User keys operations",
		Long:  `User keys operations`,
		Run: func(cmd *cobra.Command, _ []string) {
			cmd.Help()
		},
	}
	removeKeyCmd = &cobra.Command{
		Use:   "remove",
		Short: "Remove user keys",
		Long:  `Remove user keys`,
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			uid := args[0]
			accessKey := args[1]

			r := removeUserKeys(uid, accessKey)
			NewResponse(cmd, r.Success, r.Message, r.Error)
		},
	}
	addKeyCmd = &cobra.Command{
		Use:   "add",
		Short: "Add new key to user",
		Long:  `ToDo`,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			uid := args[0]

			r := addUserKey(uid)
			NewResponse(cmd, r.Success, r.Message, r.Error)
		},
	}
)

func init() {
	keysCmd.SetHelpTemplate(userKeysTemplate())
	removeKeyCmd.SetHelpTemplate(userRemoveKeyTemplate())
	removeKeyCmd.SetUsageTemplate(userRemoveKeyTemplate())
	addKeyCmd.SetHelpTemplate(userAddKeyTemplate())
	addKeyCmd.SetUsageTemplate(userAddKeyTemplate())
	userCmd.AddCommand(keysCmd)
	keysCmd.AddCommand(removeKeyCmd)
	keysCmd.AddCommand(addKeyCmd)
}

func removeUserKeys(uid string, accessKey string) CLIResponse {
	c, err := admin.New(cephHost, cephAccessKey, cephAccessSecret, nil)
	if err != nil {
		return NewResponseStruct(false, "", err.Error())
	}

	err = c.RemoveKey(context.Background(), admin.UserKeySpec{UID: uid, KeyType: "s3", AccessKey: accessKey})

	if err != nil {
		return NewResponseStruct(false, "", err.Error())
	}

	resp := fmt.Sprintf("Successfully removed keys for user %s with AccessKey %s", uid, accessKey)
	return NewResponseStruct(true, resp, "")
}

func addUserKey(uid string) CLIResponse {
	c, err := admin.New(cephHost, cephAccessKey, cephAccessSecret, nil)
	if err != nil {
		return NewResponseStruct(false, "", err.Error())
	}
	generateKey := true

	_, err = c.CreateKey(context.Background(), admin.UserKeySpec{UID: uid, KeyType: "s3", GenerateKey: &generateKey})

	if err != nil {
		return NewResponseStruct(false, "", err.Error())
	}

	r := fmt.Sprintf("Successfully added new keys for user %s", uid)
	return NewResponseStruct(true, r, "")
}
