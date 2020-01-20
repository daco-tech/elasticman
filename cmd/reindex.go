/*
Copyright © 2020 Daniel Costa <git@danielcosta.pt>
This file is part of {{ .appName }}.
*/
package cmd

import (
	"elasticman/elastic"
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

// reindexCmd represents the reindex command
var reindexCmd = &cobra.Command{
	Use:   "reindex",
	Short: "Reindex indice to another (new or add to existing index)",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Reindex called: Index: " + args[0] + " -> " + args[1])
		var result = elastic.Reindex(args[0], args[1])
		if result {
			log.Println("SUCCESS! Reindex index '" + args[0] + "' was reindexed to: '" + args[1] + "'.")
		} else {
			log.Println("FAIL! Reindex index '" + args[0] + "' was NOT reindexed to: '" + args[1] + "'.")
		}

	},
}

func init() {
	rootCmd.AddCommand(reindexCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// reindexCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// reindexCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
