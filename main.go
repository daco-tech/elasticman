package main

import (
	"elasticman/elastic"
	"elasticman/general"
	"log"
)

func main() {
	log.Println("Elastic Maintenance Tool")
	log.Println("Starting ElasticMan...")
	log.Println("--> Loading Configs...")
	config, _ := general.LoadConfiguration("config.json")
	if config.Log.Verbose {
		log.Println("*** VERBOSE LOG ACTIVATED ***")
	}

	var index = "indexname"
	var result = elastic.DeleteIndex(config.ElasticSearch.Host, index, config.Log.Verbose)
	log.Println("Delete Index: "+index+"; Result: ", result)

	var clStatus, err = elastic.GetClusterStatus(config.ElasticSearch.Host, config.Log.Verbose)
	if err == "" {
		log.Println("Cluster Status: " + clStatus.Status)
	}

}
