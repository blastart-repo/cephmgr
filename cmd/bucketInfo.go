package cmd

import (
	"context"
	"encoding/json"
	"fmt"

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
			if cmd.PersistentFlags().Changed("cluster") {
				overrideActiveCluster(clusterOverride)
			}
			listBuckets(cmd)
		},
	}
	getBucketInfoCmd = &cobra.Command{
		Use:   "info",
		Short: "Get bucket details",
		Long:  `Get bucket details`,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if cmd.PersistentFlags().Changed("cluster") {
				overrideActiveCluster(clusterOverride)
			}
			bucket := &admin.Bucket{
				Bucket: args[0],
			}

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
	listBucketsCmd.SetHelpTemplate(bucketListTemplate())
	getBucketInfoCmd.SetHelpTemplate(bucketInfoTemplate())
	listBucketsCmd.SetUsageTemplate(bucketListTemplate())
	getBucketInfoCmd.SetUsageTemplate(bucketInfoTemplate())
}

func listBuckets(cmd *cobra.Command) {
	c, err := admin.New(activeCluster.EndpointURL, activeCluster.AccessKey, activeCluster.AccessSecret, nil)
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

func getBucketInfo(cmd *cobra.Command, bucket admin.Bucket) {
	c, err := admin.New(activeCluster.EndpointURL, activeCluster.AccessKey, activeCluster.AccessSecret, nil)
	if err != nil {
		NewResponse(cmd, false, "", err.Error())
	}

	b, err := c.GetBucketInfo(context.Background(), bucket)
	if err != nil {
		NewResponse(cmd, false, "", err.Error())
		return
	}

	header := "ID\tBucket\tOwner"
	dataFormat := "%s\t%s\t%s"
	data := []interface{}{
		b.ID,
		b.Bucket,
		b.Owner,
	}

	switch {
	case returnJSON:
		bucketInfo := BucketInfo{
			ID:     b.ID,
			Bucket: b.Bucket,
			Owner:  b.Owner,
		}
		uJSON, err := json.Marshal(bucketInfo)
		if err != nil {
			NewResponse(cmd, false, "", err.Error())
		}
		fmt.Println(string(uJSON))
	default:
		printTabularData(header, dataFormat, data...)
	}
}

func getBucketInfoUsage(cmd *cobra.Command, bucket admin.Bucket) {
	c, err := admin.New(activeCluster.EndpointURL, activeCluster.AccessKey, activeCluster.AccessSecret, nil)
	if err != nil {
		NewResponse(cmd, false, "", err.Error())
	}

	b, err := c.GetBucketInfo(context.Background(), bucket)
	if err != nil {
		NewResponse(cmd, false, "", err.Error())
		return
	}

	header := "Bucket\tSize\tNumObjects"
	dataFormat := "%s\t%s\t%d"
	data := []interface{}{
		b.Bucket,
		units.BytesSize(float64(*b.Usage.RgwMain.Size)),
		int(*b.Usage.RgwMain.NumObjects),
	}

	switch {
	case returnJSON:
		bucketInfo := BucketInfoUsage{
			Bucket:     b.Bucket,
			Size:       units.BytesSize(float64(*b.Usage.RgwMain.Size)),
			NumObjects: b.Usage.RgwMain.NumObjects,
		}
		uJSON, err := json.Marshal(bucketInfo)
		if err != nil {
			NewResponse(cmd, false, "", err.Error())
		}
		fmt.Println(string(uJSON))
	default:
		printTabularData(header, dataFormat, data...)
	}
}
