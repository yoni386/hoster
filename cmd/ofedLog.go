package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	ofedCmd.AddCommand(ofedLogCmd)
}

var ofedLogCmd = &cobra.Command{
	Use:   "log",
	Short: "ofed log the host hoster",
	Long:  `All software has versions. This is hoster's`,
	Run: func(cmd *cobra.Command, args []string) {

	},
}