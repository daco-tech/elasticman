/*
Copyright © 2020 Daniel Costa <git@danielcosta.pt>
This file is part of ElasticMan.
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

var specific string

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete expired indices. Use --index to specify a single index",
	Long:  `To delete multiple expired (based on the configuration) indices.`,
	Run: func(cmd *cobra.Command, args []string) {

		if specific == "" {
			log.Println("DELETE EXPIRED INDICES")
			if ConfirmAll {
				deleteAction()
			} else {
				log.Println("Type 'Y' to continue, or 'N' to abort:")
				if general.AskForConfirmation() {
					deleteAction()
				} else {
					log.Println("Aborted!")
				}
			}
		} else {
			log.Println("DELETE SINGLE INDEX")

			if ConfirmAll {
				deleteSpecificIndex()
			} else {
				log.Println("Type 'Y' to continue, or 'N' to abort:")
				if general.AskForConfirmation() {
					deleteSpecificIndex()
				} else {
					log.Println("Aborted!")
				}
			}

		}
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().StringVar(&specific, "index", "", "specify the index name to be deleted (Default: delete all expired)")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deleteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deleteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
func deleteSpecificIndex() {
	var result = elastic.DeleteIndex(specific)
	if result {
		log.Println("Index with name '" + specific + "' deleted!")
	} else {
		log.Println("Index with name '" + specific + "' NOT deleted! Check log.")
	}
}
func deleteAction() {
	var parsedIndices, _ = elastic.GetParsedIndices()
	var totalDeleted int
	if singleton.GetConfig().Actions.Delete.Enabled {
		for _, deletions := range singleton.GetConfig().Actions.Delete.Todo {
			totalDeleted += elastic.DeleteByDays(parsedIndices, deletions.KeepDays, deletions.Logtype, deletions.Loglevel)
		}
		if totalDeleted > 0 {
			log.Println("Number of deleted indices: " + strconv.Itoa(totalDeleted))
		} else {
			log.Println("Nothing deleted!")
		}
	}
}
