package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"hoster/cluster"
)

var build string
var addFlag string
var addKernelSupport bool
var withNvmeF bool

func init() {
	ofedCmd.AddCommand(ofedInstallCmd)
	ofedInstallCmd.Flags().StringVarP(&build, "build", "B", "", "value for \"build=\" (required)")
	ofedInstallCmd.MarkFlagRequired("build")
	ofedInstallCmd.Flags().BoolVarP(&addKernelSupport, "add-kernel-support", "k", false, "add-kernel-support flag")
	ofedInstallCmd.Flags().BoolVarP(&withNvmeF, "with-nvmf", "n", false, "with-nvmf flag")
	ofedInstallCmd.Flags().StringVarP(&addFlag, "addFlag", "a", "", "add any flag e.g. --someFlag \"--\" is required")
}

var ofedInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "ofed install the host hoster",
	Long:  `All software has versions. This is hoster's`,
	Run: func(cmd *cobra.Command, args []string) {

		cluster := cluster.New(HostNames)

		// TODO: add error handle to hosts.Init() and to return?
		if err := cluster.Init(); err != nil {
			log.Errorf("hosts.Init() error: %s\n", err)
			return // TODO: return and error block doesn't reach here
		}

		cluster.OFEDInstall(build, addFlag, addKernelSupport, withNvmeF)
	},
}