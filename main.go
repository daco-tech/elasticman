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

	app.Action = func(c *cli.Context) error {
		log.Println("Elastic Maintenance Tool")
		setup()
		return nil
	}
	clierr := app.Run(os.Args)
	if clierr != nil {
		log.Fatal(clierr)
	}

}
func setup() {

	log.Println("Starting ElasticMan...")
	log.Println("--> Loading Configs...")
	config, _ := general.LoadConfiguration("config.json")
	if config.Log.Verbose {
		log.Println("*** VERBOSE LOG ACTIVATED ***")
	}

	var clStatus, err = elastic.GetClusterStatus(config.Elasticsearch.Host, config.Log.Verbose)

	if (config.Elasticsearch.RequiredStatus != "" && config.Elasticsearch.RequiredStatus != clStatus.Status) || (clStatus.NumberOfPendingTasks > config.Elasticsearch.MaxNumberOfPendingTasks) {
		log.Fatalln("Not in the desired status. Exiting Execution.")
	}

	if err == "" && config.Log.Verbose {
		log.Println("Cluster Status: " + clStatus.Status)
	}
	action(config)

}

func action(config general.Config) {

	var parsed_indices, _ = elastic.GetParsedIndices(config.Elasticsearch.Host, config.Log.Verbose, config.Parser.DateFormat, config.Parser.DateIndexLastChars, config.Parser.Loglevels, config.Parser.Logtypes)
	if config.Actions.Delete.Enabled {
		for _, deletions := range config.Actions.Delete.Todo {
			elastic.DeleteByDays(config.Elasticsearch.Host, config.Actions.Delete.DryRun, parsed_indices, deletions.KeepDays, deletions.Logtype, deletions.Loglevel, config.Log.Verbose)
		}

	}

	//var index = "indexname"
	//var result = elastic.DeleteIndex(config.ElasticSearch.Host, index, config.Log.Verbose)
	//log.Println("Delete Index: "+index+"; Result: ", result)

}
