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
		Args:  cobra.ExactArgs(1), // Require exactly 1 argument (UID)
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				NewResponse(cmd, false, "", "UID is required")
				return
			}

			quota := &QuotaSpec{
				UID: args[0],
			}
			getUserQuotas(cmd, quota)

		},
	}
	userQuotaSetCmd = &cobra.Command{
		Use:   "set",
		Short: "set user quotas",
		Long:  `Set user quotas`,
		Args:  cobra.ExactArgs(1), // Require exactly 1 argument (UID)
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				NewResponse(cmd, false, "", "UID is required")
				return
			}

			quota := &QuotaSpec{
				UID: args[0],
			}

			if cmd.Flags().Changed("max-objects") {
				quota.MaxObjects = &maxObjectsFlag
			}

			if cmd.Flags().Changed("max-size") {
				bytes, err := units.RAMInBytes(maxSizeFlag)
				if err != nil {
					NewResponse(cmd, false, "", err.Error())
					return
				}
				quota.MaxSize = &bytes
			}

			if cmd.Flags().Changed("enabled") {
				quota.Enabled = &enabledFlag
			}

			resp := setUserQuotas(quota)
			NewResponse(cmd, resp.Success, resp.Message, resp.Error)
		},
	}
)

func init() {

	userCmd.AddCommand(userQuotaCmd)
	userQuotaCmd.AddCommand(userQuotaGetCmd)
	userQuotaCmd.AddCommand(userQuotaSetCmd)
	userQuotaSetCmd.SetHelpTemplate(userQuotaSetTemplate())
	userQuotaSetCmd.SetUsageTemplate(userQuotaSetTemplate())
	userQuotaGetCmd.SetHelpTemplate(userQuotaGetTemplate())
	userQuotaGetCmd.SetUsageTemplate(userQuotaGetTemplate())
	userQuotaSetCmd.Flags().Int64Var(&maxObjectsFlag, "max-objects", -1, "Max Objects Quota. Usage: --max-objects=<int>")
	userQuotaSetCmd.Flags().StringVar(&maxSizeFlag, "max-size", "", "Max Size Quota ")
	userQuotaSetCmd.Flags().BoolVar(&enabledFlag, "enabled", false, "Enable or disable quotas")
}

func getUserQuotas(cmd *cobra.Command, quotaSpec *QuotaSpec) {
	c, err := admin.New(cephHost, cephAccessKey, cephAccessSecret, nil)
	if err != nil {
		NewResponse(cmd, false, "", err.Error())
	}

	u, err := c.GetUserQuota(context.Background(), admin.QuotaSpec{UID: quotaSpec.UID})

	if err != nil {
		NewResponse(cmd, false, "", err.Error())
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
			NewResponse(cmd, false, "", err.Error())
		}
		fmt.Println(string(uJSON))
	default:
		w := tabwriter.NewWriter(os.Stdout, 10, 1, 5, ' ', 0)
		fs := "%s\t%s\t%v\t%v\n"
		fmt.Fprintln(w, "UID\tMaxSize\tMaxObjects\tEnabled")
		fmt.Fprintf(w, fs, quotaSpec.UID, units.BytesSize(float64(*u.MaxSize)), *u.MaxObjects, *u.Enabled)
		w.Flush()
	}

}

func setUserQuotas(quotaSpec *QuotaSpec) CLIResponse {
	c, err := admin.New(cephHost, cephAccessKey, cephAccessSecret, nil)
	if err != nil {
		return NewResponseStruct(false, "", err.Error())
	}

	adminQuotaSpec := admin.QuotaSpec{
		UID:        quotaSpec.UID,
		MaxObjects: quotaSpec.MaxObjects,
		MaxSize:    quotaSpec.MaxSize,
		Enabled:    quotaSpec.Enabled,
	}

	err = c.SetUserQuota(context.Background(), adminQuotaSpec)
	if err != nil {
		return NewResponseStruct(false, "", err.Error())
	}

	successMessage := "Quota set successfully"
	return NewResponseStruct(true, successMessage, "")
}
