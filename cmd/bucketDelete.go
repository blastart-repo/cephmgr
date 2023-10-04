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
		Run:   runDeleteBucketCmd,
	}
)

func init() {
	bucketCmd.AddCommand(deleteBucketsCmd)
	deleteBucketsCmd.Flags().BoolVar(&populatedFlag, "populated", false, "Delete populated buckets")
}

func runDeleteBucketCmd(cmd *cobra.Command, args []string) {
	populated, _ := cmd.Flags().GetBool("populated")
	response := deleteBucket(args[0], populated)
	NewResponse(cmd, response.Success, response.Message, response.Error)
}

func deleteBucket(bucketName string, populated bool) CLIResponse {
	c, err := admin.New(cephHost, cephAccessKey, cephAccessSecret, nil)
	if err != nil {
		return NewCLIResponse(false, "", err.Error())
	}

	purgeObject := &populated
	bucket := admin.Bucket{
		Bucket:      bucketName,
		PurgeObject: purgeObject,
	}

	err = c.RemoveBucket(context.Background(), bucket)
	if err != nil {
		return NewCLIResponse(false, "", err.Error())
	}

	successMessage := fmt.Sprintf("Bucket '%s' deleted successfully.", bucketName)
	return NewCLIResponse(true, successMessage, "")
}
