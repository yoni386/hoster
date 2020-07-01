package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(mstCmd)
}

var mstCmd = &cobra.Command{
	Use:   "mst",
	Short: "mst the host hoster",
	Long:  `All software has versions. This is hoster's`,
	Version: "0.9",
	//Args: cobra.MinimumNArgs(1),

	//ValidArgs: []string{"start",},
	Run: func(cmd *cobra.Command, args []string) {

	},
}