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
	listBucketsCmd = &cobra.Command{
		Use:   "list",
		Short: "Get a list of buckets",
		Long:  `get list of buckets.`,
		Run: func(cmd *cobra.Command, _ []string) {
			listBuckets(cmd)

		},
	}
	getBucketInfoCmd = &cobra.Command{
		Use:   "info",
		Short: "Get bucket details",
		Long:  `Get bucket details`,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			bucket := &Bucket{
				Bucket: args[0],
			}
			/*if bucket.Bucket == "" {
				fmt.Printf("error: %s\n", errMissingBucketID)
				cmd.Help()
				os.Exit(1)
			}*/
			switch {
			case bucketUsageInfo:
				getBucketInfoUsage(cmd, *bucket)
			case bucketQuotaInfo:
				getBucketQuotas(cmd, *bucket)
			default:
				getBucketInfo(cmd, *bucket)
			}
		},
	}
)

func init() {
	bucketCmd.AddCommand(listBucketsCmd)
	bucketCmd.AddCommand(getBucketInfoCmd)
	getBucketInfoCmd.PersistentFlags().BoolVarP(&bucketUsageInfo, "usage", "u", false, "Bucket usage")
	getBucketInfoCmd.PersistentFlags().BoolVarP(&bucketQuotaInfo, "quota", "q", false, "Bucket quotas")

}

func listBuckets(cmd *cobra.Command) {
	c, err := admin.New(cephHost, cephAccessKey, cephAccessSecret, nil)
	if err != nil {
		NewResponse(cmd, false, "", err.Error())
	}
	buckets, err := c.ListBuckets(context.Background())
	if err != nil {
		NewResponse(cmd, false, "", err.Error())
	}

	switch {
	case returnJSON:

		uJSON, err := json.Marshal(&buckets)
		if err != nil {
			NewResponse(cmd, false, "", err.Error())
		}
		fmt.Println(string(uJSON))
	default:
		for _, j := range buckets {
			fmt.Println(j)
		}
	}
}

func getBucketInfo(cmd *cobra.Command, bucket Bucket) {
	c, err := admin.New(cephHost, cephAccessKey, cephAccessSecret, nil)
	if err != nil {
		NewResponse(cmd, false, "", err.Error())
	}

	b, err := c.GetBucketInfo(context.Background(), admin.Bucket{Bucket: bucket.Bucket})
	if err != nil {
		NewResponse(cmd, false, "", err.Error())
	}
	switch {
	case returnJSON:
		bucket := BucketInfo{
			ID:     b.ID,
			Bucket: b.Bucket,
			Owner:  b.Owner,
		}
		uJSON, err := json.Marshal(bucket)
		if err != nil {
			NewResponse(cmd, false, "", err.Error())
		}
		fmt.Println(string(uJSON))
	default:
		w := tabwriter.NewWriter(os.Stdout, 10, 1, 5, ' ', 0)
		fs := "%s\t%s\t%s\n"
		fmt.Fprintln(w, "ID\tBucket\tOwner")
		fmt.Fprintf(w, fs, b.ID, b.Bucket, b.Owner)
		w.Flush()
	}

}

func getBucketInfoUsage(cmd *cobra.Command, bucket Bucket) {
	c, err := admin.New(cephHost, cephAccessKey, cephAccessSecret, nil)
	if err != nil {
		NewResponse(cmd, false, "", err.Error())
	}

	b, err := c.GetBucketInfo(context.Background(), admin.Bucket{Bucket: bucket.Bucket})
	if err != nil {
		NewResponse(cmd, false, "", err.Error())
	}
	switch {
	case returnJSON:
		bucket := BucketInfoUsage{
			Bucket:     b.Bucket,
			Size:       units.BytesSize(float64(*b.Usage.RgwMain.Size)),
			NumObjects: b.Usage.RgwMain.NumObjects,
		}
		uJSON, err := json.Marshal(bucket)
		if err != nil {
			NewResponse(cmd, false, "", err.Error())
		}
		fmt.Println(string(uJSON))
	default:
		w := tabwriter.NewWriter(os.Stdout, 10, 1, 5, ' ', 0)
		fs := "%s\t%s\t%d\n"
		fmt.Fprintln(w, "Bucket\tSize\tNumObjects")
		fmt.Fprintf(w, fs, b.Bucket, units.BytesSize(float64(*b.Usage.RgwMain.Size)), int(*b.Usage.RgwMain.NumObjects))
		w.Flush()
	}
}
