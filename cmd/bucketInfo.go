package cmd

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/ceph/go-ceph/rgw/admin"
	"github.com/spf13/cobra"
)

var (
	listBucketsCmd = &cobra.Command{
		Use:   "list",
		Short: "Get a list of buckets",
		Long:  `get list of buckets.`,
		Run: func(cmd *cobra.Command, _ []string) {
			err := listBuckets()
			if err != nil {
				fmt.Println(err)
				cmd.Help()
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
			if bucket.Bucket == "" {
				fmt.Printf("error: %s\n", errMissingBucketID)
				cmd.Help()
				os.Exit(1)
			}
			switch {
			case bucketUsageInfo:
				err := getBucketInfoUsage(*bucket)
				if err != nil {
					fmt.Println(err)
					cmd.Help()
				}
			case bucketQuotaInfo:
				err := getBucketQuotas(*bucket)
				if err != nil {
					fmt.Println(err)
					cmd.Help()
				}
			default:
				err := getBucketInfo(*bucket)
				if err != nil {
					fmt.Println(err)
					cmd.Help()
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

func listBuckets() error {
	c, err := admin.New(cephHost, cephAccessKey, cephAccessSecret, nil)
	if err != nil {
		return err
	}
	buckets, err := c.ListBuckets(context.Background())

	if err != nil {
		return err
	}

	for _, j := range buckets {
		fmt.Println(j)
	}
	return nil
}

func getBucketInfo(bucket Bucket) error {
	c, err := admin.New(cephHost, cephAccessKey, cephAccessSecret, nil)
	if err != nil {
		return err
	}

	b, err := c.GetBucketInfo(context.Background(), admin.Bucket{Bucket: bucket.Bucket})
	if err != nil {
		return err
	}

	w := tabwriter.NewWriter(os.Stdout, 10, 1, 5, ' ', 0)

	fs := "%s\t%s\t%s\n"
	fmt.Fprintln(w, "ID\tBucket\tOwner")
	fmt.Fprintf(w, fs, b.ID, b.Bucket, b.Owner)
	w.Flush()

	return nil
}

func getBucketInfoUsage(bucket Bucket) error {
	c, err := admin.New(cephHost, cephAccessKey, cephAccessSecret, nil)
	if err != nil {
		return err
	}

	b, err := c.GetBucketInfo(context.Background(), admin.Bucket{Bucket: bucket.Bucket})
	if err != nil {
		return err
	}
	fmt.Printf("Bucket: %s\n", b.Bucket)
	fmt.Printf("Size : %d \n", getUint64Value(b.Usage.RgwMain.Size))
	fmt.Printf("NumObjects : %d \n", getUint64Value(b.Usage.RgwMain.NumObjects))
	fmt.Printf("SizeUtilized : %d \n", getUint64Value(b.Usage.RgwMain.SizeUtilized))
	fmt.Printf("SizeActual : %d \n", getUint64Value(b.Usage.RgwMain.SizeActual))

	return nil
}

func getUint64Value(ptr *uint64) uint64 {
	if ptr != nil {
		return *ptr
	}
	return 0
}
