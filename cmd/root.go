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
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Config struct {
	Hostname     string `mapstructure:"hostname"`
	AccessKey    string `mapstructure:"accessKey"`
	AccessSecret string `mapstructure:"accessSecret"`
}

var (
	cfgFile string
	rootCmd = &cobra.Command{
		Use:   "cephmgr",
		Short: "Ceph RGW management CLI tool",
		Long: `This tool manages Ceph cluster RGW parameters from command line.
		
To manage cluster, you must provide clusteri address and credentials. 
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
		viper.ReadInConfig()
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)
		configName := ".cephmgr.yaml"
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(configName)
		viper.AutomaticEnv()

		if err := viper.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); ok {
				fmt.Println("Creating default config file")
				hostname := ReadKey("Ceph S3 Host (with scheme):")
				viper.Set("hostname", hostname)
				accesskey := ReadKey("Access key:")
				viper.Set("accessKey", accesskey)
				accesssecret := ReadKey("Access secret:")
				viper.Set("accessSecret", accesssecret)
				err = viper.WriteConfigAs(filepath.Join(home, configName))
				if err != nil {
					fmt.Printf("Cannot write configuration file: %v\n", err)
				}
			} else {
				fmt.Fprintln(os.Stderr, "configfile exists, but something else is wrong")
			}
		}
	}

	var config Config
	err := viper.Unmarshal(&config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not decode config into struct: %v\n", err)
	}

	if err := viper.ReadInConfig(); err != nil {
		fmt.Fprintln(os.Stderr, "Cannot read config file:", viper.ConfigFileUsed())
	}
	cephHost = config.Hostname
	cephAccessKey = config.AccessKey
	cephAccessSecret = config.AccessSecret
}

func ReadKey(label string) string {
	var s string
	r := bufio.NewReader(os.Stdin)
	for {
		fmt.Fprint(os.Stderr, label+" ")
		s, _ = r.ReadString('\n')
		if s != "" {
			break
		}
	}
	return strings.TrimSpace(s)
}
