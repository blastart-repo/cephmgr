package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

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

			user := &User{
				ID:          args[0],
				DisplayName: userFullname,
				Email:       userEmail,
			}

			err := modifyUser(*user)
			if err != nil {
				fmt.Println(err)
				cmd.Help()
			}

		},
	}
)

func init() {
	userCmd.AddCommand(modifyCmd)

}

func modifyUser(user User) error {

	c, err := admin.New(cephHost, cephAccessKey, cephAccessSecret, nil)
	if err != nil {
		return err
	}
	resp, err := c.ModifyUser(context.Background(), admin.User{ID: user.ID, DisplayName: user.DisplayName, Email: user.Email})

	if err != nil {
		return err
	}

	buser, _ := json.Marshal(resp)

	var userdata User
	_ = json.Unmarshal([]byte(buser), &userdata)

	fmt.Println("User info modifyed")

	w := tabwriter.NewWriter(os.Stdout, 10, 1, 5, ' ', 0)

	fs := "%s\t%s\t%s\n"
	fmt.Fprintln(w, "UID\tFull Name\tEmail")
	fmt.Fprintf(w, fs, userdata.ID, userdata.DisplayName, userdata.Email)
	w.Flush()
	return nil
}
