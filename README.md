# ElasticMan - Elastic Maintenance Tool

This is a simple Go program to do the ElasticSearch Log files maintenance.

The objective to this project is to create a tool with enough flexibility to clean, consolidate and organize ElasticSearch indices.

## Getting started

* Make sure you have [dep](https://github.com/golang/dep) installed
* Clone this repo `git clone https://github.com/daco-tech/elasticman.git`
* Create a configuration file with name: config.json at main.go file level with the text in the config section
* Run `make` to download dependencies and run the application


## Config

config.json file content example:
```
{
    "elasticsearch": {
        "host": "http://<elasticsearch_server>:9200"
    },
    "log": {
        "verbose": true
    }
}
```