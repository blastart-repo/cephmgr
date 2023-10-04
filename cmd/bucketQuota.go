package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"

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
			bucket := &Bucket{
				Bucket: args[0],
			}

			err := getBucketQuotas(cmd, *bucket)
			if err != nil {
				NewResponse(cmd, false, "", err.Error())
			}
		},
	}
	bucketQuotaSetCmd = &cobra.Command{
		Use:   "set",
		Short: "set bucket quotas",
		Long:  `Set bucket quotas`,
		Run: func(cmd *cobra.Command, args []string) {
			quota := &QuotaSpec{
				UID:    args[0],
				Bucket: args[1],
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
}

func getBucketQuotas(cmd *cobra.Command, bucket Bucket) error {

	c, err := admin.New(cephHost, cephAccessKey, cephAccessSecret, nil)
	if err != nil {
		return err
	}

	b, err := c.GetBucketInfo(context.Background(), admin.Bucket{Bucket: bucket.Bucket})
	if err != nil {
		return err
	}
	respQuota := ResponseQuota{
		Bucket:     b.Bucket,
		Enabled:    b.BucketQuota.Enabled,
		MaxSize:    units.BytesSize(float64(*b.BucketQuota.MaxSize)),
		MaxObjects: b.BucketQuota.MaxObjects,
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
		fs := "%s\t%v\t%s\t%t\n"
		fmt.Fprintln(w, "Bucket\tMaxSize\tMaxObjects\tEnabled")
		formattedMaxObjects := strconv.FormatInt(*b.BucketQuota.MaxObjects, 10)
		fmt.Fprintf(w, fs, b.Bucket, units.BytesSize(float64(*b.BucketQuota.MaxSize)), formattedMaxObjects, *b.BucketQuota.Enabled)
		w.Flush()
	}

	return nil
}

func setBucketQuotas(quotaSpec *QuotaSpec) CLIResponse {
	c, err := admin.New(cephHost, cephAccessKey, cephAccessSecret, nil)
	if err != nil {
		return NewCLIResponse(false, "", err.Error())
	}

	adminQuotaSpec := admin.QuotaSpec{
		UID:        quotaSpec.UID,
		Bucket:     quotaSpec.Bucket,
		MaxObjects: quotaSpec.MaxObjects,
		MaxSize:    quotaSpec.MaxSize,
		Enabled:    quotaSpec.Enabled,
	}

	// Set the user quota using the admin API
	err = c.SetIndividualBucketQuota(context.Background(), adminQuotaSpec)
	if err != nil {
		return NewCLIResponse(false, "", err.Error())
	}

	successMessage := "Quota set successfully"
	return NewCLIResponse(true, successMessage, "")
}
