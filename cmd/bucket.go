/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/ceph/go-ceph/rgw/admin"
	"github.com/spf13/cobra"
)

// bucketCmd represents the bucket command
var (
	bucketCmd = &cobra.Command{
		Use:   "bucket",
		Short: "Bucket related commands",
		Long:  `Bucket related commands`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
	listBucketsCmd = &cobra.Command{
		Use:   "list",
		Short: "Get a list of buckets",
		Long:  `get list of buckets.`,
		Run: func(cmd *cobra.Command, args []string) {
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
				ID: args[0],
			}
			fmt.Println("args 0", args[0])
			if bucket.ID == "" {
				fmt.Printf("error: %s\n", errMissingBucketID)
				cmd.Help()
				os.Exit(1)
			}
			err := getBucketInfo(*bucket)
			if err != nil {
				fmt.Println(err)
				cmd.Help()
			}
		},
	}
)

func init() {
	rgwCmd.AddCommand(bucketCmd)
	bucketCmd.AddCommand(listBucketsCmd)
	bucketCmd.AddCommand(getBucketInfoCmd)
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
	fmt.Println("inside function")
	c, err := admin.New(cephHost, cephAccessKey, cephAccessSecret, nil)
	if err != nil {
		return err
	}

	fmt.Println("after connection")
	b, err := c.GetBucketInfo(context.Background(), admin.Bucket{ID: bucket.ID})
	if err != nil {
		return err
	}

	fmt.Println("after getbucketinfo")
	w := tabwriter.NewWriter(os.Stdout, 10, 1, 5, ' ', 0)

	fs := "%s\t%s\t%s\n"
	fmt.Fprintln(w, "ID\tBucket\tOwner")
	fmt.Fprintf(w, fs, b.ID, b.Bucket, b.Owner)
	w.Flush()

	return nil
}
