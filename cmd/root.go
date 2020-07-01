package cmd

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path"
	"runtime"
	"strings"
)

var Verbose bool
var Overwrite bool
var Info bool
var Debug bool
//var author string
var HostNames []string
var CfgFile string
var Prefix string

//SupportedExts "json", "toml", "yaml", "yml", "properties", "props", "prop", "hcl"
func init() {
	//cobra.OnInitialize(initConfig)
	HostNames = []string{}
	rootCmd.PersistentFlags().StringVarP(&CfgFile, "Config", "C", "", "configuration files allowed are JSON, TOML, YAML and HCL. change later TODO: config file")
	rootCmd.PersistentFlags().StringSliceVarP(&HostNames, "host", "H", HostNames, "host list")
	rootCmd.PersistentFlags().StringVarP(&Prefix, "Prefix", "P", "", "Prefix for hostname e.g. \"fqdn.local\"")
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().BoolVarP(&Overwrite, "overwrite", "o", false, "Flag config overwrite config from config file")
	rootCmd.PersistentFlags().BoolVarP(&Info, "info", "i", false, "info output")
	rootCmd.PersistentFlags().BoolVarP(&Debug, "debug", "d", false, "debug output")
	viper.BindPFlag("machines.prefix", rootCmd.PersistentFlags().Lookup("Prefix"))
	//rootCmd.PersistentFlags().StringVar(&author, "author", "YOUR NAME", "Author name for copyright attribution")
	//viper.BindPFlag("author", rootCmd.PersistentFlags().Lookup("author"))
	//viper.BindPFlag("useViper", rootCmd.PersistentFlags().Lookup("viper"))

}

func initConfig() {
	// Don't forget to read config either from cfgFile or from home directory!
	//if cfgFile != "" {
	//	// Use config file from the flag.
	//	viper.SetConfigFile(cfgFile)
	//} else {
	//	// Find home directory.
	//	home, err := homedir.Dir()
	//	if err != nil {
	//		fmt.Println(err)
	//		os.Exit(1)
	//	}
	//
	//	// Search config in home directory with name ".cobra" (without extension).
	//	viper.AddConfigPath(home)
	//	viper.SetConfigName(".cobra")


		//if Overwrite {
		//	viper.BindPFlag("machines.prefix", rootCmd.PersistentFlags().Lookup("Prefix"))
		//}
	//}

	//if err := viper.ReadInConfig(); err != nil {
	//	fmt.Println("Can't read config:", err)
	//	os.Exit(1)
	//}
}
// TODO: add machinesPrefix string, machinesNames string to func
// TODO: change global var []string define return and call params
// func makeHostNames(machinesPrefix string, machinesNames string) []string {
func makeHostNames() []string {
	// get prefix from flag or config-file
	prefix := viper.GetString("machines.prefix")

	// get hostNames from config-file
	hostNamesSourceConfig := viper.GetStringSlice("machines.names")
	// get hostNames from flag
	hostNamesSourceFlag := HostNames
	// make slice for hostNames from config-file and flag
	hostNames := make([]string, 0, len(hostNamesSourceConfig) + len(hostNamesSourceFlag))

	if len(hostNamesSourceConfig) > 0 || len(hostNamesSourceFlag) > 0 {

		if len(hostNamesSourceConfig) > 0 && len(hostNamesSourceFlag) > 0 {
			log.Debugf("host were set via Flag and Config. Merge will be done.\n")
		}

		if len(hostNamesSourceConfig) > 0 {
			log.Debugf("hostNamesSourceConfig len is: %d\n", len(hostNamesSourceConfig))

			for Index := range hostNamesSourceConfig {
				hostname := hostNamesSourceConfig[Index]
				if prefix != "" && !(strings.Contains(hostname, ".")) {
					log.Debugf("Host: \"%s\" prefix: \"%s\" will be added\n", hostname, prefix )
					hostname = fmt.Sprintf("%s.%s", hostname, prefix)
				}
				log.Debugf("Host: \"%s\" is added to hostNamesPrefix\n", hostname)
				hostNames = append(hostNames, hostname)
			}
		}

		if len(hostNamesSourceFlag) > 0 {
			log.Debugf("hostNamesSourceConfig len is: %d\n", len(hostNamesSourceFlag))

			for Index := range hostNamesSourceFlag {
				hostname := hostNamesSourceFlag[Index]
				if prefix != "" && !(strings.Contains(hostname, ".")) {
					log.Debugf("Host: \"%s\" prefix: \"%s\" will be added\n", hostname, prefix )
					hostname = fmt.Sprintf("%s.%s", hostname, prefix)
				}
				log.Debugf("Host: \"%s\" is added to hostNamesPrefix\n", hostname)
				hostNames = append(hostNames, hostname)
			}

		}

	}

	log.Infof("hostNamesPrefix len is: %d hostNamesPrefix: %q\n", len(hostNames), hostNames)
	return hostNames
}


type PlainFormatter struct {
	TimestampFormat string
	LevelDesc []string
}

type PlainFormatterDebug struct {
	TimestampFormat string
	LevelDesc []string
}

func (f *PlainFormatter) Format(entry *log.Entry) ([]byte, error) {
	timestamp := fmt.Sprintf(entry.Time.Format(f.TimestampFormat))

	return []byte(fmt.Sprintf("[%s] %s - %s", f.LevelDesc[entry.Level], timestamp, entry.Message)), nil
}

func (f *PlainFormatterDebug) Format(entry *log.Entry) ([]byte, error) {
	timestamp := fmt.Sprintf(entry.Time.Format(f.TimestampFormat))

	pc, file, line, _ := runtime.Caller(7)
	funcName := runtime.FuncForPC(pc).Name()

	entry.Data["source"] = fmt.Sprintf("%s:%v:%s", path.Base(file), line, path.Base(funcName))
	// TODO: add real pid and check if entry.Data["source"] or entry.Data can used used instead of path.Base and etc

	return []byte(fmt.Sprintf("[%s] %s %s:%d %s - %s", f.LevelDesc[entry.Level], timestamp, file, line, path.Base(funcName), entry.Message)), nil
}


//func init()  {
//
//	if Verbose {
//		plainDebugFormatter := new(PlainFormatterDebug)
//		plainDebugFormatter.TimestampFormat = "02-01-2006 15:04:05"
//		plainDebugFormatter.LevelDesc = []string{"PANC", "FATL", "ERRO", "WARN", "INFO", "DEBG"}
//		log.SetFormatter(plainDebugFormatter)
//		log.SetOutput(os.Stdout)
//	} else {
//		plainFormatter := new(PlainFormatter)
//		plainFormatter.TimestampFormat = "15:04:05"
//		plainFormatter.LevelDesc = []string{"PANC", "FATL", "ERRO", "WARN", "INFO", "DEBG"}
//		log.SetFormatter(plainFormatter)
//		log.SetOutput(os.Stdout)
//	}
//
//	//Debug = true
//
//	//fmt.Printf("debug: %v verbose: %v\n", Debug, Verbose)
//
//	if Info {
//		log.SetLevel(log.InfoLevel)
//	} else if Debug {
//		log.SetLevel(log.DebugLevel)
//	} else {
//		log.SetLevel(log.WarnLevel)
//	}
//
//}

var rootCmd = &cobra.Command{
	Use:   "hoster",
	Short: "hoster is tool to manage host ",
	Long: `hoster info on hoster long long long`,
	//Args: cobra.MinimumNArgs(1),
	PersistentPreRun: func(cmd *cobra.Command, args []string) {

		//fmt.Printf("!!!!rootCmd PersistentPreRun!!!!\n")
		if Verbose {
			plainDebugFormatter := new(PlainFormatterDebug)
			plainDebugFormatter.TimestampFormat = "02-01-2006 15:04:05"
			plainDebugFormatter.LevelDesc = []string{"PANC", "FATL", "ERRO", "WARN", "INFO", "DEBG"}
			log.SetFormatter(plainDebugFormatter)
			log.SetOutput(os.Stdout)
		} else {
			plainFormatter := new(PlainFormatter)
			plainFormatter.TimestampFormat = "15:04:05"
			plainFormatter.LevelDesc = []string{"PANC", "FATL", "ERRO", "WARN", "INFO", "DEBG"}
			log.SetFormatter(plainFormatter)
			log.SetOutput(os.Stdout)
		}

		//Debug = true

		//fmt.Printf("debug: %v verbose: %v\n", Debug, Verbose)

		if Info {
			log.SetLevel(log.InfoLevel)
		} else if Debug {
			log.SetLevel(log.DebugLevel)
		} else {
			log.SetLevel(log.WarnLevel)
		}

		//fmt.Printf("Post rootCmd PersistentPreRun debug: %v verbose: %v\n", Debug, Verbose)

		if CfgFile != "" {
			viper.SetConfigFile(CfgFile)
			err := viper.ReadInConfig() // Find and read the config file
			if err != nil { // Handle errors reading the config file
				log.Errorf("Fatal error read config file: %s\n", err)
				os.Exit(1) // TODO: add error if data1.toml is used to continue or exit?
				//return
			}
		}
		// TODO: change global var []string define return and call params
		HostNames = makeHostNames()
		//ClusterName =
		//log.Debugf("HostNames: %#v\n", HostNames)


	},
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here


		//fmt.Printf("222 debug: %v verbose: %v\n", Debug, Verbose)
		//
		//if Info {
		//	log.SetLevel(log.InfoLevel)
		//} else if Debug {
		//	log.SetLevel(log.DebugLevel)
		//} else {
		//	log.SetLevel(log.WarnLevel)
		//}
		//
		//fmt.Println(args)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
