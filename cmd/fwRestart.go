package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	fwCmd.AddCommand(fwRestartCmd)
}

var fwRestartCmd = &cobra.Command{
	Use:   "restart",
	Short: "fw restart the host hoster",
	Long:  `All software has versions. This is hoster's`,
	Run: func(cmd *cobra.Command, args []string) {

	},
}