package cmd

import (
	"context"

	"github.com/ceph/go-ceph/rgw/admin"
)

func createBucket(bucketName string) error {
	c, err := admin.New(cephHost, cephAccessKey, cephAccessSecret, nil)
	if err != nil {
		return err
	}

	// Define the bucket configuration
	bucketConfig := admin.Bucket{
		ID: bucketName, // Bucket name is used as the ID
	}

	// Create the bucket
	err = c.CreateBucket(context.Background(), bucketConfig)
	if err != nil {
		return err
	}

	return nil
}
