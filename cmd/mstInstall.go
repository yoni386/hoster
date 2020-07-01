package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"hoster/cluster"
)

func init() {
	mstCmd.AddCommand(mstInstallCmd)
}

var mstInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "mst install the host hoster",
	Long:  `All software has versions. This is hoster's`,
	Run: func(cmd *cobra.Command, args []string) {
		cluster := cluster.New(HostNames)

		// TODO: add error handle to hosts.Init() and to return?
		if err := cluster.Init(); err != nil {
			log.Errorf("hosts.Init() error: %s\n", err)
			return // TODO: return and error block doesn't reach here
		}

		cluster.MSTInstaller()
	},
}