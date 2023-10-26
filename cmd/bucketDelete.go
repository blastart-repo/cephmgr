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
			if cmd.PersistentFlags().Changed("cluster") {
				overrideActiveCluster(clusterOverride)
			}

			populated, _ := cmd.Flags().GetBool("populated")
			response := deleteBucket(args[0], populated)
			NewResponse(cmd, response.Success, response.Message, response.Error)
		},
	}
)

func init() {
	bucketCmd.AddCommand(deleteBucketsCmd)
	deleteBucketsCmd.Flags().BoolVar(&populatedFlag, "populated", false, "Delete populated buckets")
	deleteBucketsCmd.SetHelpTemplate(bucketDeleteTemplate())
	deleteBucketsCmd.SetUsageTemplate(bucketDeleteTemplate())
}

func deleteBucket(bucketName string, populated bool) CLIResponse {
	c, err := admin.New(activeCluster.EndpointURL, activeCluster.AccessKey, activeCluster.AccessSecret, nil)
	if err != nil {
		return NewResponseStruct(false, "", err.Error())
	}

	bucket := admin.Bucket{
		Bucket:      bucketName,
		PurgeObject: &populated,
	}

	err = c.RemoveBucket(context.Background(), bucket)
	if err != nil {
		return NewResponseStruct(false, "", err.Error())
	}

	successMessage := fmt.Sprintf("Bucket '%s' deleted successfully.", bucketName)
	return NewResponseStruct(true, successMessage, "")
}
