package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/ceph/go-ceph/rgw/admin"
	"github.com/docker/go-units"
	"github.com/spf13/cobra"
)

var (
	bucketQuotaCmd = &cobra.Command{
		Use:   "quota",
		Short: "Bucket quota operations",
		Long:  `Bucket quota operations`,
		Run: func(cmd *cobra.Command, _ []string) {
			cmd.Help()
		},
	}
	bucketQuotaGetCmd = &cobra.Command{
		Use:   "get",
		Short: "get bucket quotas",
		Long:  `todo`,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if cmd.PersistentFlags().Changed("cluster") {
				overrideActiveCluster(clusterOverride)
			}
			bucket := &admin.Bucket{
				Bucket: args[0],
			}
			getBucketQuotas(cmd, *bucket)
		},
	}
	bucketQuotaSetCmd = &cobra.Command{
		Use:   "set",
		Short: "set bucket quotas",
		Long:  `Set bucket quotas`,
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			if cmd.PersistentFlags().Changed("cluster") {
				overrideActiveCluster(clusterOverride)
			}
			quota := &admin.QuotaSpec{
				UID:    args[0],
				Bucket: args[1],
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

			response := setBucketQuotas(quota)
			NewResponse(cmd, response.Success, response.Message, response.Error)
		},
	}
)

func init() {
	bucketCmd.AddCommand(bucketQuotaCmd)
	bucketQuotaCmd.AddCommand(bucketQuotaGetCmd)
	bucketQuotaCmd.AddCommand(bucketQuotaSetCmd)

	bucketQuotaSetCmd.Flags().Int64Var(&maxObjectsFlag, "max-objects", -1, "Max Objects Quota. Usage: --max-objects=<int>")
	bucketQuotaSetCmd.Flags().StringVar(&maxSizeFlag, "max-size", "", "Max Size Quota (in bytes)")
	bucketQuotaSetCmd.Flags().BoolVar(&enabledFlag, "enabled", false, "Enable or disable quotas")
	bucketQuotaGetCmd.SetHelpTemplate(bucketQuotaGetTemplate())
	bucketQuotaSetCmd.SetHelpTemplate(bucketQuotaSetTemplate())
	bucketQuotaGetCmd.SetUsageTemplate(bucketQuotaGetTemplate())
	bucketQuotaSetCmd.SetUsageTemplate(bucketQuotaSetTemplate())
}

func getBucketQuotas(cmd *cobra.Command, bucket admin.Bucket) {

	c, err := admin.New(activeCluster.EndpointURL, activeCluster.AccessKey, activeCluster.AccessSecret, nil)
	if err != nil {
		NewResponse(cmd, false, "", err.Error())
	}

	b, err := c.GetBucketInfo(context.Background(), bucket)
	if err != nil {
		NewResponse(cmd, false, "", err.Error())
		return
	}

	header := "Bucket\tMaxSize\tMaxObjects\tEnabled"
	dataFormat := "%s\t%v\t%s\t%t"
	formattedMaxObjects := strconv.FormatInt(*b.BucketQuota.MaxObjects, 10)
	data := []interface{}{
		b.Bucket,
		units.BytesSize(float64(*b.BucketQuota.MaxSize)),
		formattedMaxObjects,
		*b.BucketQuota.Enabled,
	}

	switch {
	case returnJSON:
		respQuota := ResponseQuota{
			Bucket:     b.Bucket,
			Enabled:    b.BucketQuota.Enabled,
			MaxSize:    units.BytesSize(float64(*b.BucketQuota.MaxSize)),
			MaxObjects: b.BucketQuota.MaxObjects,
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

func setBucketQuotas(quotaSpec *admin.QuotaSpec) CLIResponse {
	c, err := admin.New(activeCluster.EndpointURL, activeCluster.AccessKey, activeCluster.AccessSecret, nil)
	if err != nil {
		return NewResponseStruct(false, "", err.Error())
	}
	valueKb := int(bytesToKB(*quotaSpec.MaxSize))
	quotaSpec.MaxSizeKb = &valueKb
	err = c.SetIndividualBucketQuota(context.Background(), *quotaSpec)
	if err != nil {
		return NewResponseStruct(false, "", err.Error())
	}
	successMessage := "Quota set successfully"
	return NewResponseStruct(true, successMessage, "")
}
