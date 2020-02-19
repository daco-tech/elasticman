/*
Copyright © 2020 Daniel Costa <git@danielcosta.pt>
This file is part of {{ .appName }}.
*/
package cmd

import (
	"elasticman/elastic"
	"elasticman/general"
	"elasticman/singleton"
	"log"
	"strconv"

	"github.com/spf13/cobra"
)

// consolidateCmd represents the consolidate command
var consolidateCmd = &cobra.Command{
	Use:   "consolidate",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("CONSOLIDATE INDICES")
		if ConfirmAll {
			consolidateAction()
		} else {
			log.Println("Type 'Y' to continue, or 'N' to abort:")
			if general.AskForConfirmation() {
				consolidateAction()
			} else {
				log.Println("Aborted!")
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(consolidateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// consolidateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// consolidateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func consolidateAction() {
	var parsedIndices, _ = elastic.GetParsedIndices()
	var totalConsolidated int
	if singleton.GetConfig().Actions.Consolidate.Enabled {
		for _, consolidations := range singleton.GetConfig().Actions.Consolidate.Todo {
			totalConsolidated += elastic.ConsolidateByDays(parsedIndices, consolidations.KeepDays, consolidations.Logtype, consolidations.Loglevel, consolidations.Suffix, singleton.GetConfig().Actions.Consolidate.DryRun, singleton.GetConfig().Actions.Consolidate.DeleteSourceIndex, consolidations.AddCurrentYear, consolidations.RemoveLogLevel, consolidations.RemoveDate)
		}
		if totalConsolidated > 0 {
			log.Println("Number of consolidated indices: " + strconv.Itoa(totalConsolidated))
		} else {
			log.Println("Nothing consolidated!")
		}
	}
}
