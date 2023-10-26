package cmd

import (
	"context"

	"github.com/ceph/go-ceph/rgw/admin"
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var (
	modifyCmd = &cobra.Command{
		Use:   "modify",
		Short: "Modify user",
		Long:  `Modify user`,
		Run: func(cmd *cobra.Command, args []string) {
			if cmd.PersistentFlags().Changed("cluster") {
				overrideActiveCluster(clusterOverride)
			}
			var user *admin.User
			if len(args) > 0 {
				user = &admin.User{
					ID:          args[0],
					DisplayName: userFullname,
					Email:       userEmail,
				}
			} else {
				cmd.Help()
				return
			}

			resp := modifyUser(*user)
			NewResponse(cmd, resp.Success, resp.Message, resp.Error)
		},
	}
)

func init() {
	userCmd.AddCommand(modifyCmd)
	modifyCmd.SetHelpTemplate(userModifyTemplate())

}

func modifyUser(user admin.User) CLIResponse {

	c, err := admin.New(activeCluster.EndpointURL, activeCluster.AccessKey, activeCluster.AccessSecret, nil)
	if err != nil {
		return NewResponseStruct(false, "", err.Error())
	}
	_, err = c.ModifyUser(context.Background(), user)
	if err != nil {
		return NewResponseStruct(false, "", err.Error())
	}

	return NewResponseStruct(true, "User info modified", "")
}
