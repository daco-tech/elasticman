/*
Copyright © 2020 DANIEL COSTA <git@danielcosta.pt>
*/
package cmd

import (
	"elasticman/general"
	"elasticman/singleton"
	"fmt"
	"log"
	"os"
	"os/user"

	"github.com/spf13/cobra"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ElasticMan",
	Short: "Elastic Maintenance Tool",
	Long: `This is a tool to maintain ElasticSearch easily.
The name is based on the Elastic (witch is a trade mark from Elastic) name, and the Portuguese
word "manutenção" that translated to english means "maintenance" it is also an analogy to 
the man who maintains the elastic search node or cluster.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.elasticman/config.json)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("verbose", "v", false, "Sets verbose output to true (overrides configuration)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile == "" {
		usr, userErr := user.Current()
		if userErr != nil {
			log.Fatal(userErr)
		}
		cfgFile = usr.HomeDir + "/.elasticman/config.json"
	}
	log.Println("--> Loading Configs from " + cfgFile + "...")
	config, cfgErr := general.LoadConfiguration(cfgFile)
	if cfgErr != nil {
		log.Fatalln("No configuration file found. Looking for '" + cfgFile + "'.")
	}
	singleton.SetConfig(config)
}
