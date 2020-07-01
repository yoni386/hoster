package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(ofedCmd)
}

var ofedCmd = &cobra.Command{
	Use:   "ofed",
	Short: "ofed the host hoster",
	Long:  `All software has versions. This is hoster's`,
	Run: func(cmd *cobra.Command, args []string) {

	},
}