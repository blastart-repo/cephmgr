package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
)

var (
	clusterCmd = &cobra.Command{
		Use:   "cluster",
		Short: "Ceph clusters operations",
		Long:  `Get clusters information. Change default active cluster.`,
		Run: func(cmd *cobra.Command, _ []string) {
			cmd.Help()
		},
	}

	listClustersCmd = &cobra.Command{
		Use:   "list",
		Short: "Get a list of clusters",
		Long:  `Get a list of available clusters.`,
		Run: func(cmd *cobra.Command, _ []string) {
			listClusters(cmd)
		},
	}

	getActiveClusterCmd = &cobra.Command{
		Use:   "get_active",
		Short: "Get default active cluster info",
		Long:  `Get default active cluster info`,
		Run: func(cmd *cobra.Command, _ []string) {
			getActiveCluster(cmd)
		},
	}

	setActiveClusterCmd = &cobra.Command{
		Use:   "set_active",
		Short: "Set default active cluster",
		Long:  `Set default active cluster`,
		Run: func(cmd *cobra.Command, _ []string) {
			setActiveCluster(cmd, clusterName)
		},
	}

	addClusterCmd = &cobra.Command{
		Use:   "add",
		Short: "Add new cluster",
		Long:  `Add new cluster to the list of available clusters`,
		Run: func(cmd *cobra.Command, _ []string) {
			addNewCluster(cmd, Cluster{
				ClusterName:  clusterName,
				AccessKey:    clusterAccessKey,
				AccessSecret: clusterAccessSecret,
				EndpointURL:  clusterEndpointURL,
			})
		},
	}

	removeClusterCmd = &cobra.Command{
		Use:   "remove <cluster name>",
		Short: "Removes cluster",
		Long:  `Removes the cluster from the list of available clusters`,
		Run: func(cmd *cobra.Command, _ []string) {
			removeCluster(cmd, clusterName)
		},
	}
)

func init() {
	rgwCmd.AddCommand(clusterCmd)
	clusterCmd.AddCommand(getActiveClusterCmd)
	clusterCmd.AddCommand(listClustersCmd)
	clusterCmd.AddCommand(setActiveClusterCmd)
	clusterCmd.AddCommand(addClusterCmd)
	clusterCmd.AddCommand(removeClusterCmd)

	listClustersCmd.SetHelpTemplate(clusterListTemplate())
	listClustersCmd.SetUsageTemplate(clusterListTemplate())
	getActiveClusterCmd.SetHelpTemplate(clusterGetActiveTemplate())
	getActiveClusterCmd.SetUsageTemplate(clusterGetActiveTemplate())
	setActiveClusterCmd.SetHelpTemplate(clusterSetActiveTemplate())
	setActiveClusterCmd.SetUsageTemplate(clusterSetActiveTemplate())
	addClusterCmd.SetHelpTemplate(clusterAddNewTemplate())
	addClusterCmd.SetUsageTemplate(clusterAddNewTemplate())
	removeClusterCmd.SetHelpTemplate(clusterRemoveTemplate())
	removeClusterCmd.SetUsageTemplate(clusterRemoveTemplate())

	setActiveClusterCmd.Flags().StringVarP(&clusterName, "name", "n", "", "Cluster name")
	setActiveClusterCmd.MarkFlagRequired("name")

	addClusterCmd.Flags().StringVarP(&clusterName, "name", "n", "", "Cluster name")
	addClusterCmd.MarkFlagRequired("name")
	addClusterCmd.Flags().StringVarP(&clusterAccessKey, "access_key", "k", "", "Cluster access key")
	addClusterCmd.MarkFlagRequired("access_key")
	addClusterCmd.Flags().StringVarP(&clusterAccessSecret, "access_secret", "s", "", "Cluster access secret")
	addClusterCmd.MarkFlagRequired("access_secret")
	addClusterCmd.Flags().StringVarP(&clusterEndpointURL, "endpoint_url", "e", "", "Cluster endpoint URL")
	addClusterCmd.MarkFlagRequired("endpoint_url")

	removeClusterCmd.Flags().StringVarP(&clusterName, "name", "n", "", "Cluster name")
	removeClusterCmd.MarkFlagRequired("name")
}

func listClusters(cmd *cobra.Command) { // TODO don't show access keys and secrets
	if returnJSON {
		cJSON, err := json.Marshal(clusterConfig.Clusters)
		if err != nil {
			NewResponse(cmd, false, "", err.Error())
		}
		fmt.Println(string(cJSON))
	} else {
		fmt.Println(fmt.Sprintf("%+v", clusterConfig.Clusters))
	}
}

func getActiveCluster(cmd *cobra.Command) {
	if returnJSON {
		cJSON, err := json.Marshal(activeCluster)
		if err != nil {
			NewResponse(cmd, false, "", err.Error())
		}
		fmt.Println(string(cJSON))
	} else {
		fmt.Println(fmt.Sprintf("%+v", activeCluster))
	}
}

func setActiveCluster(cmd *cobra.Command, name string) {
	changeActiveCluster(name)
	if returnJSON {
		cJSON, err := json.Marshal(activeCluster)
		if err != nil {
			NewResponse(cmd, false, "", err.Error())
		}
		fmt.Println(string(cJSON))
	} else {
		fmt.Println(fmt.Sprintf("%+v", activeCluster))
	}
}

// addNewCluster - uses the helper function newCluster to add the new cluster to the config file
func addNewCluster(cmd *cobra.Command, cluster Cluster) {
	newCluster(cluster)
	if returnJSON {
		cJSON, err := json.Marshal(clusterConfig.Clusters)
		if err != nil {
			NewResponse(cmd, false, "", err.Error())
		}
		fmt.Println(string(cJSON))
	} else {
		fmt.Println(fmt.Sprintf("%+v", clusterConfig.Clusters))
	}
}

// removeCluster - uses the helper function remCluster to change the config file to remove the cluster from the list of available clusters
func removeCluster(cmd *cobra.Command, name string) {
	remCluster(name)
	if returnJSON {
		cJSON, err := json.Marshal(clusterConfig.Clusters)
		if err != nil {
			NewResponse(cmd, false, "", err.Error())
		}
		fmt.Println(string(cJSON))
	} else {
		fmt.Println(fmt.Sprintf("%+v", clusterConfig.Clusters))
	}
}
