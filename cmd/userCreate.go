package cmd

import (
	"context"
	"fmt"

	"github.com/ceph/go-ceph/rgw/admin"
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var (
	createCmd = &cobra.Command{
		Use:   "create",
		Short: "Create new user",
		Long:  `Create new user.`,
		Run: func(cmd *cobra.Command, _ []string) {
			if cmd.PersistentFlags().Changed("cluster") {
				overrideActiveCluster(clusterOverride)
			}

			user := &User{
				ID:          userID,
				DisplayName: userFullname,
				Email:       userEmail,
				UserCaps:    userCaps,
			}

			resp := createUser(*user)
			NewResponse(cmd, resp.Success, resp.Message, resp.Error)

		},
	}
)

func init() {
	userCmd.AddCommand(createCmd)
	createCmd.PersistentFlags().StringVarP(&userID, "user", "u", "", "Ceph user ID")
	createCmd.MarkFlagRequired("user")
	createCmd.SetHelpTemplate(userCreateTemplate())
}

func createUser(user User) CLIResponse {

	c, err := admin.New(activeCluster.EndpointURL, activeCluster.AccessKey, activeCluster.AccessSecret, nil)
	if err != nil {
		return NewResponseStruct(false, "", err.Error())
	}
	_, err = c.CreateUser(context.Background(), admin.User{ID: user.ID, DisplayName: user.DisplayName, UserCaps: user.UserCaps})

	if err != nil {
		return NewResponseStruct(false, "", err.Error())
	}

	return NewResponseStruct(true, fmt.Sprintf("Created user for %s", user.DisplayName), "")
}
