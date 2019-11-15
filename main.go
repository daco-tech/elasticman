package main

import (
	"elasticman/elastic"
	"elasticman/general"
	"log"
	"os"
	"os/user"

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
			Name:  "verbose",
			Usage: "Verbose Mode (Accepts true/false). Overide configuration verbose setting.",
		},
		cli.StringFlag{
			Name:  "delete, d",
			Value: "no",
			Usage: "'--delete yes' to delete multiple indexes (configuration). '--delete indexname' to delete a single index (dry_run does not work with this option).",
		},
	}

	app.Action = func(c *cli.Context) error {

		log.SetFlags(0)
		log.Println("   ____ __           __   _       __  ___          ")
		log.Println("  / __// /___ _ ___ / /_ (_)____ /  |/  /___ _ ___ ")
		log.Println(" / _/ / // _ `/(_-</ __// // __// /|_/ // _ \\`/ _ \\")
		log.Println("/___//_/ \\_,_//___/\\__//_/ \\__//_/  /_/ \\_,_//_//_/")
		log.SetFlags(1)
		log.Println("Starting ElasticMan...")

		usr, userErr := user.Current()
		if userErr != nil {
			log.Fatal(userErr)
		}

		//Config
		configFile := usr.HomeDir + "/.elasticman/config.json"
		if c.NArg() > 0 {
			configFile = c.Args()[0]
		}

		if c.String("config") != "" {
			configFile = c.String("config")
		}
		log.Println("--> Loading Configs...")

		config, cfgErr := general.LoadConfiguration(configFile)
		if cfgErr != nil {
			log.Fatalln("No configuration file found. Looking for '" + configFile + "'.")
		}
		//Set verbose mode
		var verbose bool = config.Log.Verbose
		if c.String("verbose") != "" && c.String("verbose") == "true" {
			verbose = true
		} else if c.String("verbose") != "" && c.String("verbose") == "false" {
			verbose = false
		}
		if verbose {
			log.Println("*** VERBOSE LOG ACTIVATED ***")
		}

		//Set run pre-conditions
		var clStatus, err = elastic.GetClusterStatus(config.Elasticsearch.Host, verbose)

		if (config.Elasticsearch.RequiredStatus != "" && config.Elasticsearch.RequiredStatus != clStatus.Status) || (clStatus.NumberOfPendingTasks > config.Elasticsearch.MaxNumberOfPendingTasks) {
			log.Fatalln("Not in the desired status. Exiting Execution.")
		}

		if err == "" && verbose {
			log.Println("Cluster Status: " + clStatus.Status)
		}
		var doSomething bool
		//Run Actions
		if c.String("delete") != "" && c.String("delete") != "no" {
			if c.String("delete") == "yes" {
				log.Println("DELETE INDICES MODE ACTIVATED")
				doSomething = true
				deleteAction(config, verbose)
			} else {
				log.Println("DELETE SINGLE INDEX MODE ACTIVATED")
				doSomething = true
				var result = elastic.DeleteIndex(config.Elasticsearch.Host, c.String("delete"), verbose)
				if result {
					log.Println("Index with name '" + c.String("delete") + "' deleted!")
				} else {
					log.Println("Index with name '" + c.String("delete") + "' NOT deleted! Check log.")
				}

			}

		}

		//If nothing runs say something
		if !doSomething {
			log.Fatalln("No action selected for execution! Nothing to do!")
		}

		return nil
	}
	cliErr := app.Run(os.Args)
	if cliErr != nil {
		log.Fatal(cliErr)
	}

}

func deleteAction(config general.Config, verbose bool) {
	var parsedIndices, _ = elastic.GetParsedIndices(config.Elasticsearch.Host, verbose, config.Parser.DateFormat, config.Parser.DateIndexLastChars, config.Parser.Loglevels, config.Parser.Logtypes)
	if config.Actions.Delete.Enabled {
		for _, deletions := range config.Actions.Delete.Todo {
			elastic.DeleteByDays(config.Elasticsearch.Host, config.Actions.Delete.DryRun, parsedIndices, deletions.KeepDays, deletions.Logtype, deletions.Loglevel, config.Log.Verbose)
		}

	}
}
