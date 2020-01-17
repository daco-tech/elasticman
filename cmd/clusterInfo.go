/*
Copyright © 2020 Daniel Costa <git@danielcosta.pt>
This file is part of {{ .appName }}.
*/
package cmd

import (
	"elasticman/elastic"
	"elasticman/singleton"
	"log"
	"strconv"

	"github.com/brettski/go-termtables"
	"github.com/spf13/cobra"
)

// clusterInfoCmd represents the clusterInfo command
var clusterInfoCmd = &cobra.Command{
	Use:   "clusterInfo",
	Short: "Returns connected cluster info.",
	Long:  `Returns the cluster details in a terminal table.`,
	Run: func(cmd *cobra.Command, args []string) {
		PrintClusterInfo()
	},
}

func init() {
	rootCmd.AddCommand(clusterInfoCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// clusterInfoCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// clusterInfoCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func PrintClusterInfo() {
	var clStatus, err = elastic.GetClusterStatus(singleton.GetConfig().Elasticsearch.Host, singleton.GetConfig().Log.Verbose)
	if err != "" {
		log.Panicln(err)
	}

	table := termtables.CreateTable()
	table.AddHeaders("CLUSTER INFO", "")
	table.AddRow("Cluster Name:", clStatus.ClusterName)
	table.AddRow("Cluster Status:", clStatus.Status)
	table.AddRow("Cluster TimedOut:", clStatus.TimedOut)
	table.AddRow("Cluster Number Of Nodes:", clStatus.NumberOfNodes)
	table.AddRow("Cluster Number Of Data Nodes:", clStatus.NumberOfDataNodes)
	table.AddRow("Cluster Active Primary Shards:", clStatus.ActivePrimaryShards)
	table.AddRow("Cluster Active Shards:", clStatus.ActiveShards)
	table.AddRow("Cluster Relocating Shards:", clStatus.RelocatingShards)
	table.AddRow("Cluster Initializing Shards:", clStatus.InitializingShards)
	table.AddRow("Cluster Unassigned Shards:", clStatus.UnassignedShards)
	table.AddRow("Cluster Delayed Unassigned Shards:", clStatus.DelayedUnassignedShards)
	table.AddRow("Cluster Number Of Pending Tasks:", clStatus.NumberOfPendingTasks)
	table.AddRow("Cluster Number Of In Flight Fetch:", clStatus.NumberOfInFlightFetch)
	table.AddRow("Cluster Task Max Waiting In Queue:", strconv.Itoa(clStatus.TaskMaxWaitingInQueueMillis)+" (Millis)")
	table.AddRow("Cluster Active Shards:", strconv.FormatFloat(clStatus.ActiveShardsPercentAsNumber, 'f', 2, 64)+" %")

	log.SetFlags(0)
	log.Println(table.Render())
	log.SetFlags(1)
}
