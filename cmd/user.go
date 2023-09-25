package cmd

import "github.com/spf13/cobra"

var (
	userCmd = &cobra.Command{
		Use:   "user",
		Short: "Ceph users operations",
		Long:  `Get users information. Create new users. Change users caps`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
)

func init() {
	rgwCmd.AddCommand(userCmd)

	userCmd.PersistentFlags().StringVarP(&userCaps, "caps", "", "", "User capabilities")
	userCmd.PersistentFlags().StringVarP(&userFullname, "fullname", "f", "", "Ceph user full name")
	userCmd.PersistentFlags().StringVarP(&userEmail, "email", "e", "", "Ceph user e-mail")

}
