package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"hoster/cluster"
)

func init() {
	rootCmd.AddCommand(coldCmd)
}

var coldCmd = &cobra.Command{
	Use:   "cold",
	Short: "cold power",
	Long:  `All software has versions. This is hoster's`,
	Run: func(cmd *cobra.Command, args []string) {

		cluster := cluster.New(HostNames)

		// TODO: add error handle to hosts.Init() and to return?
		if err := cluster.InitNonSsh(); err != nil {
			log.Errorf("hosts.InitNonSsh() error: %s\n", err)
			return // TODO: return and error block doesn't reach here
		}

		cluster.HostRestart()
	},
}