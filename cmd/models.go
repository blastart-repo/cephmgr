package cmd

import "github.com/ceph/go-ceph/rgw/admin"

type ResponseQuota struct {
	UID        string `json:"user_id" url:"uid"`
	Bucket     string `json:"bucket" url:"bucket"`
	Enabled    *bool  `json:"enabled" url:"enabled"`
	MaxSize    string `json:"max_size" url:"max-size"`
	MaxObjects *int64 `json:"max_objects" url:"max-objects"`
}

type CLIResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

type BucketInfo struct {
	ID     string `json:"id"`
	Bucket string `json:"bucket" url:"bucket"`
	Owner  string `json:"owner"`
}

type BucketInfoUsage struct {
	Bucket     string  `json:"bucket" url:"bucket"`
	Size       string  `json:"size"`
	NumObjects *uint64 `json:"num_objects"`
}

type UserCapsResponse struct {
	UID  string              `json:"user_id" url:"uid"`
	Caps []admin.UserCapSpec `json:"caps"`
}

type UserInfoResponse struct {
	UID         string              `json:"user_id" url:"uid"`
	DisplayName string              `json:"display_name" url:"display-name"`
	Email       string              `json:"email" url:"email"`
	Caps        []admin.UserCapSpec `json:"caps"`
}

type Cluster struct {
	ClusterName  string `mapstructure:"clusterName" json:"cluster_name"`
	AccessKey    string `mapstructure:"accessKey" json:"access_key"`
	AccessSecret string `mapstructure:"accessSecret" json:"access_secret"`
	EndpointURL  string `mapstructure:"endpointURL" json:"endpoint_url"`
}

type SensitiveCluster struct {
	ClusterName string `json:"cluster_name"`
	EndpointURL string `json:"endpoint_url"`
}
