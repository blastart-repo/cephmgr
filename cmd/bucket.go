/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// bucketCmd represents the bucket command
var (
	bucketCmd = &cobra.Command{
		Use:   "bucket",
		Short: "Bucket related commands",
		Long:  `Bucket related commands`,
		Run: func(cmd *cobra.Command, _ []string) {
			cmd.Help()
		},
	}
)

func init() {
	rgwCmd.AddCommand(bucketCmd)

}
