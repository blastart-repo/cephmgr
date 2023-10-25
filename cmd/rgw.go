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
	"github.com/spf13/cobra"
)

// rgwCmd represents the rgw command
var rgwCmd = &cobra.Command{
	Use:   "rgw",
	Short: "rgw module",
	Long:  `A Ceph rgw module`,
	Run: func(cmd *cobra.Command, _ []string) {
		cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(rgwCmd)
	rgwCmd.PersistentFlags().BoolVarP(&returnJSON, "json", "j", false, "Return values as json")
	rgwCmd.PersistentFlags().StringVarP(&clusterOverride, "cluster", "c", "", "Default cluster override, will use the defined cluster instead of the default one")
	// rgwCmd.PersistentFlags().StringVarP(&cephHost, "ceph", "c", "", "Ceph host name with scheme")
	// rgwCmd.PersistentFlags().StringVarP(&cephAccessKey, "accesskey", "k", "", "Ceph access key")
	// rgwCmd.PersistentFlags().StringVarP(&cephAccessSecret, "accesssecret", "s", "", "Ceph access secret")

}
