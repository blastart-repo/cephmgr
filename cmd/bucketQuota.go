package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/ceph/go-ceph/rgw/admin"
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
			if bucket.Bucket == "" {
				fmt.Printf("error: %s\n", errMissingBucketID)
				cmd.Help()
				os.Exit(1)
			}
			err := getBucketQuotas(*bucket)
			if err != nil {
				fmt.Println(err)
				cmd.Help()
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
				quota.MaxSize = &maxSizeFlag
			}

			if cmd.Flags().Changed("enabled") {
				quota.Enabled = &enabledFlag
			}

			if cmd.Flags().Changed("max-size-kb") {
				quota.MaxSizeKb = &maxSizeKbFlag
			}

			err := setBucketQuotas(quota)
			if err != nil {
				fmt.Println(err)
				cmd.Help()
			}
		},
	}
)

func init() {
	bucketCmd.AddCommand(bucketQuotaCmd)
	bucketQuotaCmd.AddCommand(bucketQuotaGetCmd)
	bucketQuotaCmd.AddCommand(bucketQuotaSetCmd)

	// Add flags to bucketQuotaSetCmd
	bucketQuotaSetCmd.Flags().Int64Var(&maxObjectsFlag, "max-objects", -1, "Max Objects Quota. Usage: --max-objects=<int>")
	bucketQuotaSetCmd.Flags().Int64Var(&maxSizeFlag, "max-size", -1, "Max Size Quota (in bytes)")
	bucketQuotaSetCmd.Flags().IntVar(&maxSizeKbFlag, "max-size-kb", 0, "Max Size KB Quota")
	bucketQuotaSetCmd.Flags().BoolVar(&enabledFlag, "enabled", false, "Enable or disable quotas")
}

func getBucketQuotas(bucket Bucket) error {

	c, err := admin.New(cephHost, cephAccessKey, cephAccessSecret, nil)
	if err != nil {
		return err
	}

	b, err := c.GetBucketInfo(context.Background(), admin.Bucket{Bucket: bucket.Bucket})
	if err != nil {
		return err
	}

	fmt.Printf("Bucket: %s\n", b.Bucket)
	fmt.Printf("Max Size : %d B\n", *b.BucketQuota.MaxSize)
	fmt.Printf("Max Objects : %d\n", *b.BucketQuota.MaxObjects)
	fmt.Printf("Max Size KB : %d KB\n", *b.BucketQuota.MaxSize)
	fmt.Printf("Enabled: %t\n", *b.BucketQuota.Enabled)

	return nil
}

func setBucketQuotas(quotaSpec *QuotaSpec) error {
	c, err := admin.New(cephHost, cephAccessKey, cephAccessSecret, nil)
	if err != nil {
		return err
	}

	// Create an admin.QuotaSpec and populate it with values from quotaSpec
	adminQuotaSpec := admin.QuotaSpec{
		UID:        quotaSpec.UID,
		Bucket:     quotaSpec.Bucket,
		MaxObjects: quotaSpec.MaxObjects,
		MaxSize:    quotaSpec.MaxSize,
		Enabled:    quotaSpec.Enabled,
		MaxSizeKb:  quotaSpec.MaxSizeKb,
	}

	// Set the user quota using the admin API
	err = c.SetIndividualBucketQuota(context.Background(), adminQuotaSpec)
	if err != nil {
		return err
	}

	fmt.Printf("User: %s\n", quotaSpec.UID)
	fmt.Println("Quota set successfully.")
	return nil
}
