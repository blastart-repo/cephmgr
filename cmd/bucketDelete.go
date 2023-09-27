package cmd

import (
	"context"
	"fmt"

	"github.com/ceph/go-ceph/rgw/admin"
	"github.com/spf13/cobra"
)

var (
	deleteBucketsCmd = &cobra.Command{
		Use:   "delete",
		Short: "Delete empty buckets",
		Long:  "Delete an empty bucket by name",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			populated, _ := cmd.Flags().GetBool("populated")
			err := deleteBucket(args[0], populated)
			if err != nil {
				fmt.Println(err)
				cmd.Help()
			}
		},
	}
	populatedFlag bool
)

func init() {
	bucketCmd.AddCommand(deleteBucketsCmd)
	deleteBucketsCmd.Flags().BoolVar(&populatedFlag, "populated", false, "Delete populated buckets")
}

func deleteBucket(bucketName string, populated bool) error {
	c, err := admin.New(cephHost, cephAccessKey, cephAccessSecret, nil)
	if err != nil {
		return err
	}

	// Create a pointer to a bool and set its value based on the populated flag
	purgeObject := &populated

	// Create an admin.Bucket with PurgeObject set to the pointer to the bool
	bucket := admin.Bucket{
		Bucket:      bucketName,
		PurgeObject: purgeObject,
	}

	err = c.RemoveBucket(context.Background(), bucket)
	if err != nil {
		return err
	}

	fmt.Printf("Bucket '%s' deleted successfully.\n", bucketName)
	return nil
}
