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

var (
	getUserCmd = &cobra.Command{
		Use:   "get",
		Short: "Get user info",
		Long:  `Get user info`,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if cmd.PersistentFlags().Changed("cluster") {
				overrideActiveCluster(clusterOverride)
			}
			user := &admin.User{
				ID: args[0],
			}
			getUser(cmd, *user)
		},
	}
	listCmd = &cobra.Command{
		Use:   "list",
		Short: "Get a list of users",
		Long:  `get list of users from the cluster.`,
		Run: func(cmd *cobra.Command, _ []string) {
			if cmd.PersistentFlags().Changed("cluster") {
				overrideActiveCluster(clusterOverride)
			}
			listUsers(cmd)

		},
	}
)

func init() {
	userCmd.AddCommand(getUserCmd)
	getUserCmd.SetHelpTemplate(getUserHelpTemplate())
	getUserCmd.SetUsageTemplate(getUserHelpTemplate())
	userCmd.AddCommand(listCmd)
	listCmd.SetHelpTemplate(listUsersTemplate())
}

func getUser(cmd *cobra.Command, user admin.User) {
	c, err := admin.New(activeCluster.EndpointURL, activeCluster.AccessKey, activeCluster.AccessSecret, nil)
	if err != nil {
		NewResponse(cmd, false, "", err.Error())
	}

	u, err := c.GetUser(context.Background(), user)
	if err != nil {
		NewResponse(cmd, false, "", err.Error())
		return
	}

	header := "UID\tFull Name\tEmail\tCaps"
	dataFormat := "%s\t%s\t%s\t%v"
	data := []interface{}{u.ID, u.DisplayName, u.Email, u.Caps}

	switch {
	case returnJSON:
		respStruct := UserInfoResponse{
			UID:         u.ID,
			DisplayName: u.DisplayName,
			Email:       u.Email,
			Caps:        u.Caps,
		}
		uJSON, err := json.Marshal(respStruct)
		if err != nil {
			NewResponse(cmd, false, "", err.Error())
		}
		fmt.Println(string(uJSON))
	default:
		printTabularData(header, dataFormat, data...)
	}
}

func listUsers(cmd *cobra.Command) {
	c, err := admin.New(activeCluster.EndpointURL, activeCluster.AccessKey, activeCluster.AccessSecret, nil)
	if err != nil {
		NewResponse(cmd, false, "", err.Error())
	}

	users, err := c.GetUsers(context.Background())
	if err != nil {
		NewResponse(cmd, false, "", err.Error())
	}
	switch {
	case returnJSON:

		uJSON, err := json.Marshal(users)
		if err != nil {
			NewResponse(cmd, false, "", err.Error())
		}
		fmt.Println(string(uJSON))
	default:
		for _, j := range *users {
			fmt.Println(j)
		}
	}
}
