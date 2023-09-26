/*
Copyright © 2022 Tarmo Katmuk <tarmo.katmuk@gmail.com>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	Version string
	Date    string
	Commit  string
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show app version",
	Long:  `Show cephmgr application version`,
	Run: func(_ *cobra.Command, _ []string) {
		fmt.Printf("cephmgr %s %s (%s)\n", Version, Commit[:7], Date)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
