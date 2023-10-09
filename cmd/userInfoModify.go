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
			var user *User
			if len(args) > 0 {
				user = &User{
					ID:          args[0],
					DisplayName: userFullname,
					Email:       userEmail,
				}
			} else {
				cmd.Help() // Show custom help when no argument is provided
				return
			}

			// Your modification logic here

			resp := modifyUser(*user)
			NewResponse(cmd, resp.Success, resp.Message, resp.Error)
		},
	}
)

func init() {
	userCmd.AddCommand(modifyCmd)
	modifyCmd.SetHelpTemplate(userModifyTemplate())

}

func modifyUser(user User) CLIResponse {

	c, err := admin.New(cephHost, cephAccessKey, cephAccessSecret, nil)
	if err != nil {
		return NewResponseStruct(false, "", err.Error())
	}
	_, err = c.ModifyUser(context.Background(), admin.User{ID: user.ID, DisplayName: user.DisplayName, Email: user.Email})
	if err != nil {
		return NewResponseStruct(false, "", err.Error())
	}

	return NewResponseStruct(true, "User info modifyed", "")
}
