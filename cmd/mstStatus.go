package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	mstCmd.AddCommand(mstStatusCmd)
}

var mstStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "mst status the host hoster",
	Long:  `All software has versions. This is hoster's`,
	Run: func(cmd *cobra.Command, args []string) {

	},
}