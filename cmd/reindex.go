/*
Copyright © 2020 Daniel Costa <git@danielcosta.pt>
This file is part of {{ .appName }}.
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// reindexCmd represents the reindex command
var reindexCmd = &cobra.Command{
	Use:   "reindex",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("reindex called")
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
