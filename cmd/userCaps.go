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
		Long: `Get user caps
		`,
		Args: cobra.ExactArgs(1), // Require exactly 1 argument (UID)
		Run: func(cmd *cobra.Command, args []string) {
			user := &User{
				ID: args[0], // Use the first argument as the UID

			}
			err := getUserCaps(*user)
			if err != nil {
				fmt.Println(err)
				cmd.Help()
			}
		},
	}
	addCapsCmd = &cobra.Command{
		Use:   "add",
		Short: "Add user capabilities",
		Long: `Add user capabilities in form 
	
	"users|buckets=*|read|write|read,write"
	
	Add multiple capabilities to user:
	
	--caps "buckets=*;users=read"`,
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			user := &User{
				ID:       args[0],
				UserCaps: userCaps,
			}
			if user.ID == "" || user.UserCaps == "" {
				cmd.Help()
				os.Exit(1)
			}

			err := addUserCaps(*user)
			if err != nil {
				fmt.Println(err)
				cmd.Help()
			} else {
				fmt.Println("New user capability added.")
			}
		},
	}
	removeCapsCmd = &cobra.Command{
		Use:   "remove",
		Short: "Remove user capabilities",
		Long: `Remove user capabilities in form 
	
	"users|buckets=*|read|write|read,write"
	
	Remove multiple capabilities to user:
	
	--caps "buckets=*;users=read"`,
		Args: cobra.ExactArgs(1), // Require exactly 1 argument (UID)
		Run: func(cmd *cobra.Command, args []string) {
			user := &User{
				ID:       args[0], // Use the first argument as the UID
				UserCaps: userCaps,
			}

			if user.ID == "" || user.UserCaps == "" {
				cmd.Help()
				os.Exit(1)
			}

			err := removeUserCaps(*user)
			if err != nil {
				fmt.Println(err)
				cmd.Help()
			}
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
}

func addUserCaps(user User) error {
	c, err := admin.New(cephHost, cephAccessKey, cephAccessSecret, nil)
	if err != nil {
		return err
	}

	_, err = c.AddUserCap(context.Background(), user.ID, user.UserCaps)

	if err != nil {
		return err
	}

	return nil
}

func removeUserCaps(user User) error {
	c, err := admin.New(cephHost, cephAccessKey, cephAccessSecret, nil)
	if err != nil {
		return err
	}

	userCaps, err := c.RemoveUserCap(context.Background(), user.ID, user.UserCaps)

	if err != nil {
		return err
	}

	fmt.Printf("User ID: %s\n", user.ID)
	fmt.Println(userCaps)
	return nil
}
func getUserCaps(user User) error {
	c, err := admin.New(cephHost, cephAccessKey, cephAccessSecret, nil)
	if err != nil {
		return err
	}

	u, err := c.GetUser(context.Background(), admin.User{ID: user.ID})

	if err != nil {
		return err
	}

	w := tabwriter.NewWriter(os.Stdout, 10, 1, 5, ' ', 0)
	fs := "%s\t%s\n"
	fmt.Fprintln(w, "UID\tCaps")
	fmt.Fprintf(w, fs, u.ID, u.Caps)
	w.Flush()
	return nil
}
