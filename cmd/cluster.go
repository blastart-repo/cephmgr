package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"text/tabwriter"
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

	clusterCmd.PersistentFlags().BoolVar(&showSensitive, "sensitive", false, "Show sensitive data (access keys and secrets)")

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

func listClusters(cmd *cobra.Command) {
	var s string
	var err error
	if returnJSON {
		s, err = jsonClusters(clusterConfig.Clusters)
		if err != nil {
			NewResponse(cmd, false, "", err.Error())
			return
		}
		fmt.Println(s)
	} else {
		printClusters(clusterConfig.Clusters)
	}
}

func getActiveCluster(cmd *cobra.Command) {
	var s string
	var err error
	if returnJSON {
		s, err = jsonClusters([]Cluster{activeCluster})
		if err != nil {
			NewResponse(cmd, false, "", err.Error())
			return
		}
		fmt.Println(s)
	} else {
		printClusters([]Cluster{activeCluster})
	}
}

func setActiveCluster(cmd *cobra.Command, name string) {
	var s string
	err := changeActiveCluster(name)
	if err != nil {
		NewResponse(cmd, false, "", err.Error())
		return
	}
	if returnJSON {
		s, err = jsonClusters([]Cluster{activeCluster})
		if err != nil {
			NewResponse(cmd, false, "", err.Error())
			return
		}
		fmt.Println(s)
	} else {
		printClusters([]Cluster{activeCluster})
	}
}

// addNewCluster - uses the helper function newCluster to add the new cluster to the config file
func addNewCluster(cmd *cobra.Command, cluster Cluster) {
	var s string
	err := newCluster(cluster)
	for _, c := range clusterConfig.Clusters {
		if c.ClusterName == cluster.ClusterName {
			NewResponse(cmd, false, "", fmt.Sprintf("a cluster with the name %s already exists, please choose another name", cluster.ClusterName))
			return
		}
	}
	if err != nil {
		NewResponse(cmd, false, "", err.Error())
		return
	}
	if returnJSON {
		s, err = jsonClusters([]Cluster{cluster})
		if err != nil {
			NewResponse(cmd, false, "", err.Error())
			return
		}
		fmt.Println(s)
	} else {
		printClusters([]Cluster{cluster})
	}
}

// removeCluster - uses the helper function remCluster to change the config file to remove the cluster from the list of available clusters
func removeCluster(cmd *cobra.Command, name string) {
	var s string
	err := remCluster(name)
	if err != nil {
		NewResponse(cmd, false, "", err.Error())
		return
	}
	if returnJSON {
		s, err = jsonClusters(clusterConfig.Clusters)
		if err != nil {
			NewResponse(cmd, false, "", err.Error())
			return
		}
		fmt.Println(s)
	} else {
		printClusters(clusterConfig.Clusters)
	}
}

func printClusters(clusters []Cluster) {
	w := tabwriter.NewWriter(os.Stdout, 10, 1, 5, ' ', 0)
	if showSensitive {
		fmt.Fprint(w, "Cluster name\tAccess Key\tAccess Secret\tEndpoint URL\n")
		for _, cluster := range clusters {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", cluster.ClusterName, cluster.AccessKey, cluster.AccessSecret, cluster.EndpointURL)
		}
	} else {
		fmt.Fprint(w, "Cluster name\tEndpoint URL\n")
		for _, cluster := range clusters {
			fmt.Fprintf(w, "%s\t%s\n", cluster.ClusterName, cluster.EndpointURL)
		}
	}
	w.Flush()
}

func jsonClusters(clusters []Cluster) (string, error) {
	if showSensitive {
		cJSON, err := json.Marshal(clusters)
		if err != nil {
			return "", err
		}
		return string(cJSON), err
	} else {
		var scluster []SensitiveCluster
		for _, cluster := range clusters {
			scluster = append(scluster, SensitiveCluster{
				ClusterName: cluster.ClusterName,
				EndpointURL: cluster.EndpointURL,
			})
		}
		cJSON, err := json.Marshal(scluster)
		if err != nil {
			return "", err
		}
		return string(cJSON), nil
	}
}
