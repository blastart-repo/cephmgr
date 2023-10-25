package cmd

var (
	activeCluster       Cluster
	clusterConfig       ClusterConfig
	clusterName         string
	clusterAccessKey    string
	clusterAccessSecret string
	clusterEndpointURL  string
	clusterOverride     string
	userCaps            string
	userEmail           string
	userFullname        string
	userID              string
	maxObjectsFlag      int64
	maxSizeFlag         string
	enabledFlag         bool
	bucketUsageInfo     bool
	bucketQuotaInfo     bool
	returnJSON          bool
	populatedFlag       bool
)
