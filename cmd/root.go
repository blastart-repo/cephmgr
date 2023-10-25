/*
Copyright © 2022 Tarmo Katmuk <tarmo.katmuk@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type ClusterConfig struct {
	ActiveClusterName string    `mapstructure:"activeClusterName" json:"active_cluster_name"`
	Clusters          []Cluster `mapstructure:"clusters" json:"clusters"`
}

var (
	cfgFile string

	rootCmd = &cobra.Command{
		Use:   "cephmgr",
		Short: "Ceph RGW management CLI tool",
		Long: `This tool manages Ceph cluster RGW parameters from command line.
		
To manage cluster, you must provide cluster address and credentials. 
You can create credentials with following command from Ceph node:

radosgw-admin user create --uid admin --display name "Administrator" --caps "buckets=*;users=*;usage=read;metadata=read;zone=read"

The command returns the JSON file, from where you can use access_key and secret_key for authentication.`,
	}
)

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	viper.SetEnvPrefix("CEPH")

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cephmgr.yaml)")

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
		viper.SetConfigType("yaml")
		err := viper.ReadInConfig()
		if err != nil {
			fmt.Fprintln(os.Stderr, fmt.Sprintf("Could not read configuration file: %s", err.Error()))
		}
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)
		configName := ".cephmgr.yaml"
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(configName)
		viper.AutomaticEnv()

		err = viper.ReadInConfig()
		if err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); ok {
				fmt.Println("Creating default config file")
				clustername := ReadKey("Name for the cluster: ")
				endpointurl := ReadKey("Ceph S3 Host endpoint (with scheme): ")
				accesskey := ReadKey("Access key: ")
				accesssecret := ReadKey("Access secret: ")
				viper.Set("activeClusterName", clustername)
				viper.Set("clusters", []Cluster{{
					ClusterName:  clustername,
					AccessKey:    accesskey,
					AccessSecret: accesssecret,
					EndpointURL:  endpointurl,
				}})
				err = viper.WriteConfigAs(filepath.Join(home, configName))
				if err != nil {
					fmt.Printf("Cannot write configuration file: %v\n", err)
				}
			} else {
				fmt.Fprintln(os.Stderr, "Configfile exists, but something else is wrong")
			}
		}
	}

	err := viper.Unmarshal(&clusterConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not decode config into struct: %v\n", err)
	}

	checkClusters(clusterConfig)
}

func ReadKey(label string) string {
	var s string
	var err error
	r := bufio.NewReader(os.Stdin)
	for {
		fmt.Fprint(os.Stderr, label)
		s, err = r.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, "Could not read entered value properly.")
		}
		if s != "" {
			break
		}
	}
	return strings.TrimSpace(s)
}

// checkClusters - checks if the user inputted active cluster is valid (if it exists in the list of available clusters) and sets the value for global variables to access the information from anywhere.
func checkClusters(config ClusterConfig) {
	if len(config.Clusters) == 0 {
		fmt.Fprintln(os.Stderr, "Could not find any clusters in your config file.")
	} else {
		if config.ActiveClusterName == "" {
			fmt.Fprintln(os.Stderr, "Default active cluster not defined. To set a default cluster, run \"cephmgr rgw cluster default set <cluster>\".")
		} else {
			for _, c := range config.Clusters {
				if c.ClusterName == config.ActiveClusterName {
					activeCluster = c
					return
				}
			}
			fmt.Fprintln(os.Stderr, "Default active cluster does not exist in the list of clusters. Please set a new default active cluster.") // TODO return list of clusters
			os.Exit(1)
		}
	}
}

// overrideActiveCluster - won't change the active cluster in the configuration, will only affect one command.
func overrideActiveCluster(name string) {
	if activeCluster.ClusterName == name || name == "" {
		return
	}
	for _, c := range clusterConfig.Clusters {
		if c.ClusterName == name {
			activeCluster = c
			return
		}
	}
}

func changeActiveCluster(name string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return errors.New(fmt.Sprintf("could not locate home directory: %s", err.Error()))
	}
	configName := ".cephmgr.yaml"
	for _, c := range clusterConfig.Clusters {
		if c.ClusterName == name {
			activeCluster = c
			viper.Set("activeClusterName", name)
			err = viper.WriteConfigAs(filepath.Join(home, configName))
			if err != nil {
				return errors.New(fmt.Sprintf("could not write to config file: %s", err.Error()))
			}
			return nil
		}
	}
	return errors.New(fmt.Sprintf("no clusters with the name %s were found in the available clusters list", name))
}

func newCluster(cluster Cluster) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return errors.New(fmt.Sprintf("could not locate home directory: %s", err.Error()))
	}
	configName := ".cephmgr.yaml"

	clusterConfig.Clusters = append(clusterConfig.Clusters, cluster)
	viper.Set("clusters", clusterConfig.Clusters)

	err = viper.WriteConfigAs(filepath.Join(home, configName))
	if err != nil {
		return errors.New(fmt.Sprintf("could not write to config file: %s", err.Error()))
	}
	return nil
}

// remCluster - helper function for the removeCluster function that uses viper to change the configuration file to remove the cluster.
func remCluster(name string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return errors.New(fmt.Sprintf("could not locate home directory: %s", err.Error()))
	}
	configName := ".cephmgr.yaml"

	var i int
	for i = 0; i < len(clusterConfig.Clusters); i++ {
		if clusterConfig.Clusters[i].ClusterName == name {
			break
		}
	}
	clusters := append(clusterConfig.Clusters[:i], clusterConfig.Clusters[i+1:]...)
	viper.Set("clusters", clusters)

	err = viper.WriteConfigAs(filepath.Join(home, configName))
	if err != nil {
		return errors.New(fmt.Sprintf("could not write to config file: %s", err.Error()))
	}
	return nil
}
