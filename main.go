package main

import (
	"elasticman/elastic"
	"elasticman/general"
	"log"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "ElasticMan"
	app.Usage = "Elastic Maintenance Tool"
	app.Version = "3,1415926535897932384626433" //yes, it Pi :)

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config, c",
			Usage: "Load configuration from `FILE`",
		},
		cli.StringFlag{
			Name:  "delete",
			Value: "no",
			Usage: "'--delete yes' to activate the Delete Option. '--delete indexname' to delete a single index.",
		},
	}

	app.Action = func(c *cli.Context) error {
		log.Println("Elastic Maintenance Tool")
		log.Println("Starting ElasticMan...")
		configFile := "config.json"
		if c.NArg() > 0 {
			configFile = c.Args()[0]
		}
		//Config
		if c.String("config") != "" {
			configFile = c.String("config")
		}
		log.Println("--> Loading Configs...")
		config, cfgErr := general.LoadConfiguration(configFile)
		if cfgErr != nil {
			log.Fatalln("No configuration file found. Looking for '" + configFile + "'.")
		}
		//Set verbose mode
		if config.Log.Verbose {
			log.Println("*** VERBOSE LOG ACTIVATED ***")
		}

		//Set run pre-conditions
		var clStatus, err = elastic.GetClusterStatus(config.Elasticsearch.Host, config.Log.Verbose)

		if (config.Elasticsearch.RequiredStatus != "" && config.Elasticsearch.RequiredStatus != clStatus.Status) || (clStatus.NumberOfPendingTasks > config.Elasticsearch.MaxNumberOfPendingTasks) {
			log.Fatalln("Not in the desired status. Exiting Execution.")
		}

		if err == "" && config.Log.Verbose {
			log.Println("Cluster Status: " + clStatus.Status)
		}
		var do_something bool = false
		//Run Actions
		if c.String("delete") != "" {
			if c.String("delete") == "yes" {
				log.Println("DELETE INDICES MODE ACTIVATED")
				do_something = true
				deleteAction(config)
			} else {
				log.Println("DELETE SINGLE INDEX MODE ACTIVATED")
				do_something = true
				var result = elastic.DeleteIndex(config.Elasticsearch.Host, c.String("delete"), config.Log.Verbose)
				if result {
					log.Println("Index with name '" + c.String("delete") + "' deleted!")
				} else {
					log.Println("Index with name '" + c.String("delete") + "' NOT deleted! Check log.")
				}

			}

		}

		//If nothing runs say something
		if !do_something {
			log.Fatalln("No action selected for execution! Nothing to do!")
		}

		return nil
	}
	clierr := app.Run(os.Args)
	if clierr != nil {
		log.Fatal(clierr)
	}

}

func deleteAction(config general.Config) {
	var parsed_indices, _ = elastic.GetParsedIndices(config.Elasticsearch.Host, config.Log.Verbose, config.Parser.DateFormat, config.Parser.DateIndexLastChars, config.Parser.Loglevels, config.Parser.Logtypes)
	if config.Actions.Delete.Enabled {
		for _, deletions := range config.Actions.Delete.Todo {
			elastic.DeleteByDays(config.Elasticsearch.Host, config.Actions.Delete.DryRun, parsed_indices, deletions.KeepDays, deletions.Logtype, deletions.Loglevel, config.Log.Verbose)
		}

	}
}
