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
	"os"

	"github.com/ceph/go-ceph/rgw/admin"
	"github.com/spf13/cobra"
)

// capsCmd represents the caps command
var (
	capsCmd = &cobra.Command{
		Use:   "caps",
		Short: "User Capabilities operations",
		Long:  `User Capabilities operations`,
		Run: func(cmd *cobra.Command, _ []string) {
			cmd.Help()
		},
	}
	getCapsCmd = &cobra.Command{
		Use:   "get",
		Short: "Get user caps",
		Long:  `Get user caps`,
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if cmd.PersistentFlags().Changed("cluster") {
				overrideActiveCluster(clusterOverride)
			}
			var userID string
			if len(args) > 0 {
				userID = args[0]
			}
			user := &admin.User{
				ID: userID,
			}
			getUserCaps(cmd, *user)
		},
	}
	addCapsCmd = &cobra.Command{
		Use:   "add",
		Short: "Add user capabilities",
		Long:  `Add user capabilities `,
		Run: func(cmd *cobra.Command, args []string) {
			if cmd.PersistentFlags().Changed("cluster") {
				overrideActiveCluster(clusterOverride)
			}
			var user admin.User
			if len(args) > 0 {
				user = admin.User{
					ID:       args[0],
					UserCaps: userCaps,
				}
			}

			if user.ID == "" || user.UserCaps == "" {
				cmd.Help()
				os.Exit(1)
			}

			resp := addUserCaps(user)
			NewResponse(cmd, resp.Success, resp.Message, resp.Error)
		},
	}
	removeCapsCmd = &cobra.Command{
		Use:   "remove",
		Short: "Remove user capabilities",
		Long:  `Remove user capabilities`,
		Run: func(cmd *cobra.Command, args []string) {
			if cmd.PersistentFlags().Changed("cluster") {
				overrideActiveCluster(clusterOverride)
			}
			var user admin.User
			if len(args) > 0 {
				user = admin.User{
					ID:       args[0],
					UserCaps: userCaps,
				}
			}

			if user.ID == "" || user.UserCaps == "" {
				cmd.Help()
				os.Exit(1)
			}

			resp := removeUserCaps(user)
			NewResponse(cmd, resp.Success, resp.Message, resp.Error)
		},
	}
)

func init() {
	userCmd.AddCommand(capsCmd)
	capsCmd.AddCommand(addCapsCmd)
	capsCmd.AddCommand(removeCapsCmd)
	capsCmd.AddCommand(getCapsCmd)
	userCmd.MarkFlagRequired("user")
	addCapsCmd.MarkFlagRequired("caps")
	removeCapsCmd.MarkFlagRequired("caps")
	getCapsCmd.SetHelpTemplate(userGetCapsTemplate())
	addCapsCmd.SetHelpTemplate(userAddCapsTemplate())
	removeCapsCmd.SetHelpTemplate(userRemoveCapsTemplate())
}

func addUserCaps(user admin.User) CLIResponse {
	c, err := admin.New(activeCluster.EndpointURL, activeCluster.AccessKey, activeCluster.AccessSecret, nil)
	if err != nil {
		return NewResponseStruct(false, "", err.Error())
	}

	_, err = c.AddUserCap(context.Background(), user.ID, user.UserCaps)

	if err != nil {
		return NewResponseStruct(false, "", err.Error())
	}

	return NewResponseStruct(true, "New user capability added.", "")
}

func removeUserCaps(user admin.User) CLIResponse {
	c, err := admin.New(activeCluster.EndpointURL, activeCluster.AccessKey, activeCluster.AccessSecret, nil)
	if err != nil {
		return NewResponseStruct(false, "", err.Error())
	}

	userCaps, err := c.RemoveUserCap(context.Background(), user.ID, user.UserCaps)

	if err != nil {
		return NewResponseStruct(false, "", err.Error())
	}

	res := fmt.Sprintf("User ID: %s capabilitys removed. Remaining caps: %s", user.ID, userCaps)

	return NewResponseStruct(true, res, "")
}
func getUserCaps(cmd *cobra.Command, user admin.User) {
	c, err := admin.New(activeCluster.EndpointURL, activeCluster.AccessKey, activeCluster.AccessSecret, nil)
	if err != nil {
		NewResponse(cmd, false, "", err.Error())
	}

	u, err := c.GetUser(context.Background(), user)
	if err != nil {
		NewResponse(cmd, false, "", err.Error())
		return
	}

	header := "UID\tCaps"
	dataFormat := "%s\t%s"
	data := []interface{}{u.ID, u.Caps}

	switch {
	case returnJSON:
		uJSON, err := json.Marshal(u.Caps)
		if err != nil {
			NewResponse(cmd, false, "", err.Error())
		}
		fmt.Println(string(uJSON))
	default:
		printTabularData(header, dataFormat, data...)
	}
}
