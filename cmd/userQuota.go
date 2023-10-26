package cmd

import (
	"context"
	"encoding/json"
	"fmt"

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
			if cmd.PersistentFlags().Changed("cluster") {
				overrideActiveCluster(clusterOverride)
			}
			if len(args) < 1 {
				NewResponse(cmd, false, "", "UID is required")
				return
			}

			quota := &admin.QuotaSpec{
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
			if cmd.PersistentFlags().Changed("cluster") {
				overrideActiveCluster(clusterOverride)
			}
			if len(args) < 1 {
				NewResponse(cmd, false, "", "UID is required")
				return
			}

			quota := &admin.QuotaSpec{
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

func getUserQuotas(cmd *cobra.Command, quotaSpec *admin.QuotaSpec) {
	c, err := admin.New(activeCluster.EndpointURL, activeCluster.AccessKey, activeCluster.AccessSecret, nil)
	if err != nil {
		NewResponse(cmd, false, "", err.Error())
	}

	u, err := c.GetUserQuota(context.Background(), *quotaSpec)
	if err != nil {
		NewResponse(cmd, false, "", err.Error())
		return
	}

	header := "UID\tMaxSize\tMaxObjects\tEnabled"
	dataFormat := "%s\t%v\t%v\t%v"
	data := []interface{}{quotaSpec.UID, units.BytesSize(float64(*u.MaxSize)), *u.MaxObjects, *u.Enabled}

	switch {
	case returnJSON:
		respQuota := ResponseQuota{
			UID:        quotaSpec.UID,
			MaxSize:    units.BytesSize(float64(*u.MaxSize)),
			MaxObjects: u.MaxObjects,
			Enabled:    u.Enabled,
		}
		uJSON, err := json.Marshal(respQuota)
		if err != nil {
			NewResponse(cmd, false, "", err.Error())
		}
		fmt.Println(string(uJSON))
	default:
		printTabularData(header, dataFormat, data...)
	}
}

func setUserQuotas(quotaSpec *admin.QuotaSpec) CLIResponse {
	c, err := admin.New(activeCluster.EndpointURL, activeCluster.AccessKey, activeCluster.AccessSecret, nil)
	if err != nil {
		return NewResponseStruct(false, "", err.Error())
	}

	err = c.SetUserQuota(context.Background(), *quotaSpec)
	if err != nil {
		return NewResponseStruct(false, "", err.Error())
	}

	successMessage := "Quota set successfully"
	return NewResponseStruct(true, successMessage, "")
}
