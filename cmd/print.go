package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"hoster/cluster"
)

func init() {
	rootCmd.AddCommand(printCmd)
}

var printCmd = &cobra.Command{
	Use:   "print",
	Short: "Print the host info hoster",
	Long:  `All software has versions. This is hoster's`,
	Run: func(cmd *cobra.Command, args []string) {

		cluster := cluster.New(HostNames)

		// TODO: add error handle to hosts.Init() and to return?
		if err := cluster.Init(); err != nil {
			log.Errorf("cluster.Init() error: %s\n", err)
			return // TODO: return and error block doesn't reach here
		}
		// TODO: decide if InitHCA ig name exper for ninit hca and mst
		if err := cluster.InitHCA(); err != nil {
			log.Errorf("cluster.InitHCA() error: %s\n", err)
			return // TODO: return and error block doesn't reach here
		}

		cluster.PrintHostHCATable()
	},
}