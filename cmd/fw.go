package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(fwCmd)
}

var fwCmd = &cobra.Command{
	Use:   "fw",
	Short: "fw hoster",
	Long:  `All software has versions. This is hoster's`,
	Run: func(cmd *cobra.Command, args []string) {

	},
}