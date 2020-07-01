package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	fwCmd.AddCommand(fwBurnCmd)
}

var fwBurnCmd = &cobra.Command{
	Use:   "burn",
	Short: "fw burn the host hoster",
	Long:  `All software has versions. This is hoster's`,
	Run: func(cmd *cobra.Command, args []string) {

	},
}