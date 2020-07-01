package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	fwCmd.AddCommand(fwLogCmd)
}

var fwLogCmd = &cobra.Command{
	Use:   "log",
	Short: "fw log the host hoster",
	Long:  `All software has versions. This is hoster's`,
	Run: func(cmd *cobra.Command, args []string) {

	},
}