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
			err := listBuckets(cmd)
			if err != nil {
				NewResponse(cmd, false, "", err.Error())
			}
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
				err := getBucketInfoUsage(cmd, *bucket)
				if err != nil {
					NewResponse(cmd, false, "", err.Error())
				}
			case bucketQuotaInfo:
				err := getBucketQuotas(cmd, *bucket)
				if err != nil {
					NewResponse(cmd, false, "", err.Error())
				}
			default:
				err := getBucketInfo(cmd, *bucket)
				if err != nil {
					NewResponse(cmd, false, "", err.Error())
				}
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

func listBuckets(cmd *cobra.Command) error {
	c, err := admin.New(cephHost, cephAccessKey, cephAccessSecret, nil)
	if err != nil {
		return err
	}
	buckets, err := c.ListBuckets(context.Background())

	if err != nil {
		return err
	}

	switch {
	case returnJSON:
		stringObject := StringSlice{
			Buckets: buckets,
		}
		uJSON, err := json.Marshal(stringObject)
		if err != nil {
			NewResponse(cmd, false, "", err.Error())
		}
		fmt.Println(string(uJSON))
	default:
		for _, j := range buckets {
			fmt.Println(j)
		}
	}
	return nil
}

func getBucketInfo(cmd *cobra.Command, bucket Bucket) error {
	c, err := admin.New(cephHost, cephAccessKey, cephAccessSecret, nil)
	if err != nil {
		return err
	}

	b, err := c.GetBucketInfo(context.Background(), admin.Bucket{Bucket: bucket.Bucket})
	if err != nil {
		return err
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
	return nil
}

func getBucketInfoUsage(cmd *cobra.Command, bucket Bucket) error {
	c, err := admin.New(cephHost, cephAccessKey, cephAccessSecret, nil)
	if err != nil {
		return err
	}

	b, err := c.GetBucketInfo(context.Background(), admin.Bucket{Bucket: bucket.Bucket})
	if err != nil {
		return err
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

	return nil
}
