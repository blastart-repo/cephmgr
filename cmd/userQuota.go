package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/ceph/go-ceph/rgw/admin"
	"github.com/docker/go-units"
	"github.com/spf13/cobra"
)

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
				bytes, err := units.RAMInBytes(maxSizeFlag)
				if err != nil {
					fmt.Printf("Error parsing %s: %v\n", maxSizeFlag, err)
				}
				quota.MaxSize = &bytes
			}

			if cmd.Flags().Changed("enabled") {
				quota.Enabled = &enabledFlag
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

	userQuotaSetCmd.Flags().Int64Var(&maxObjectsFlag, "max-objects", -1, "Max Objects Quota. Usage: --max-objects=<int>")
	userQuotaSetCmd.Flags().StringVar(&maxSizeFlag, "max-size", "", "Max Size Quota ")
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
	respQuota := ResponseQuota{
		UID:        quotaSpec.UID,
		Bucket:     u.Bucket,
		Enabled:    u.Enabled,
		MaxSize:    units.BytesSize(float64(*u.MaxSize)),
		MaxObjects: u.MaxObjects,
	}

	switch {
	case returnJSON:
		uJSON, err := json.Marshal(respQuota)
		if err != nil {
			return err
		}
		fmt.Println(string(uJSON))
	default:
		w := tabwriter.NewWriter(os.Stdout, 10, 1, 5, ' ', 0)
		fs := "%s\t%s\t%v\t%v\n"
		fmt.Fprintln(w, "UID\tMaxSize\tMaxObjects\tEnabled")
		fmt.Fprintf(w, fs, quotaSpec.UID, units.BytesSize(float64(*u.MaxSize)), *u.MaxObjects, *u.Enabled)
		w.Flush()
	}
	return nil
}

func setUserQuotas(quotaSpec *QuotaSpec) error {
	c, err := admin.New(cephHost, cephAccessKey, cephAccessSecret, nil)
	if err != nil {
		return err
	}

	adminQuotaSpec := admin.QuotaSpec{
		UID:        quotaSpec.UID,
		MaxObjects: quotaSpec.MaxObjects,
		MaxSize:    quotaSpec.MaxSize,
		Enabled:    quotaSpec.Enabled,
	}

	err = c.SetUserQuota(context.Background(), adminQuotaSpec)
	if err != nil {
		return err
	}

	fmt.Println("Quota set successfully")
	return nil
}
