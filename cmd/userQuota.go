package cmd

import (
	"context"
	"fmt"

	"github.com/ceph/go-ceph/rgw/admin"
	"github.com/spf13/cobra"
)

// capsCmd represents the caps command
var (
	userQuotaCmd = &cobra.Command{
		Use:   "quota",
		Short: "User quota operations",
		Long:  `User quota operations`,
		Run: func(cmd *cobra.Command, _ []string) {
			cmd.Help()
		},
	}
	userQuotaGetCmd = &cobra.Command{
		Use:   "get",
		Short: "get user quotas",
		Long:  `todo`,
		Run: func(cmd *cobra.Command, args []string) {
			quota := &QuotaSpec{
				UID: args[0],
			}
			err := getUserQuotas(quota)
			if err != nil {
				fmt.Println(err)
				cmd.Help()
			}
		},
	}
	userQuotaSetCmd = &cobra.Command{
		Use:   "set",
		Short: "set user quotas",
		Long:  `Set user quotas`,
		Run: func(cmd *cobra.Command, args []string) {
			quota := &QuotaSpec{
				UID: args[0],
			}

			if cmd.Flags().Changed("max-objects") {
				quota.MaxObjects = &maxObjectsFlag
			}

			if cmd.Flags().Changed("max-size") {
				quota.MaxSize = &maxSizeFlag
			}

			if cmd.Flags().Changed("enabled") {
				quota.Enabled = &enabledFlag
			}

			if cmd.Flags().Changed("max-size-kb") {
				quota.MaxSizeKb = &maxSizeKbFlag
			}

			err := setUserQuotas(quota)
			if err != nil {
				fmt.Println(err)
				cmd.Help()
			}
		},
	}
)

func init() {
	userCmd.AddCommand(userQuotaCmd)
	userQuotaCmd.AddCommand(userQuotaGetCmd)
	userQuotaCmd.AddCommand(userQuotaSetCmd)

	// Add flags to userQuotaSetCmd
	userQuotaSetCmd.Flags().Int64Var(&maxObjectsFlag, "max-objects", -1, "Max Objects Quota. Usage: --max-objects=<int>")
	userQuotaSetCmd.Flags().Int64Var(&maxSizeFlag, "max-size", -1, "Max Size Quota (in bytes)")
	userQuotaSetCmd.Flags().IntVar(&maxSizeKbFlag, "max-size-kb", 0, "Max Size KB Quota")
	userQuotaSetCmd.Flags().BoolVar(&enabledFlag, "enabled", false, "Enable or disable quotas")
}

func getUserQuotas(quotaSpec *QuotaSpec) error {

	c, err := admin.New(cephHost, cephAccessKey, cephAccessSecret, nil)
	if err != nil {
		return err
	}

	u, err := c.GetUserQuota(context.Background(), admin.QuotaSpec{UID: quotaSpec.UID})

	if err != nil {
		return err
	}
	fmt.Printf("User: %s\n", quotaSpec.UID)
	fmt.Printf("Max Size : %d B\n", *u.MaxSize)
	fmt.Printf("Max Objects : %d\n", *u.MaxObjects)
	fmt.Printf("Max Size KB : %d KB\n", *u.MaxSizeKb)
	fmt.Printf("Enabled: %t\n", *u.Enabled)
	return nil
}

func setUserQuotas(quotaSpec *QuotaSpec) error {
	c, err := admin.New(cephHost, cephAccessKey, cephAccessSecret, nil)
	if err != nil {
		return err
	}

	// Create an admin.QuotaSpec and populate it with values from quotaSpec
	adminQuotaSpec := admin.QuotaSpec{
		UID:        quotaSpec.UID,
		MaxObjects: quotaSpec.MaxObjects,
		MaxSize:    quotaSpec.MaxSize,
		Enabled:    quotaSpec.Enabled,
		MaxSizeKb:  quotaSpec.MaxSizeKb,
	}

	// Set the user quota using the admin API
	err = c.SetUserQuota(context.Background(), adminQuotaSpec)
	if err != nil {
		return err
	}

	fmt.Printf("User: %s\n", quotaSpec.UID)
	fmt.Println("Quota set successfully.")
	return nil
}
