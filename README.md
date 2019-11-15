# ElasticMan - Elastic Maintenance Tool

This is a simple tool to maintain the ElasticSearch Log indices.

The main objective of this project is to create a tool with enough flexibility to clean, consolidate and organize ElasticSearch indices easily.

This tool does not compete with ElasticSearch Curator that is a more complete and maintained tool, but it is a tool that you can keep on your laptop or even in a docker container with very low hardware resources usage.

## Getting started

### From Source
    * Make sure you have [dep](https://github.com/golang/dep) installed
    * Clone this repo `git clone https://github.com/daco-tech/elasticman.git`
    * Create a configuration file with name: config.json at main.go file level with the text in the config section
    * Run `make` to download dependencies and run the application

## Config

config.json file content example:
```
    {
        "elasticsearch": {
            "host": "http://<server>:9200",
            "required_status": "green",
            "max_number_of_pending_tasks": 0
        },
        "log": {
            "verbose": true
        },
        "parser": {
            "date_format": "YYYY.MM.DD",
            "date_index_last_chars": 10,
            "loglevels": [
                "debug",
                "trace",
                "verbose",
                "info",
                "warn",
                "warning",
                "error",
                "critical"
            ],
            "logtypes": [
                "log",
                "httpstatus",
                "metric"
            ]
        },
        "actions": {
            "delete": {
                "enabled": true,
                "dry_run": true,
                "todo": [
                    {
                        "loglevel": "debug",
                        "keep-days": 2
                    },
                    {
                        "loglevel": "trace",
                        "keep-days": 2
                    },
                    {
                        "loglevel": "verbose",
                        "keep-days": 2
                    },
                    {
                        "loglevel": "info",
                        "keep-days": 7,
                        "logtype": "log"
                    },
                    {
                        "loglevel": "warn",
                        "keep-days": 14
                    },
                    {
                        "loglevel": "warning",
                        "keep-days": 14
                    },
                    {
                        "loglevel": "critical",
                        "keep-days": 60
                    },
                    {
                        "loglevel": "error",
                        "keep-days": 60
                    }
                ]
            }
        }
    }
```

This configuration works with indexes with name like "app-example-log-warn-2013.09.20" where the date_index_last_chars are the last 10 chars of this index.

## Actions

### DELETE 

#### DELETE MULTIPLE INDICES
Usage: elasticman --delete yes

This command deletes indices according with the configuration defined in the actions -> delete section. 

This tool can be tested without delete indices if the  "dry_run" option is set to true in the configuration file.

Once you test the configuration file, you are able to schedule in a cron job the automatic execution of this cleanup tool to keep your elastic search indices clean (ideal for logs).

Be aware that the indices date should be in the end of their name and should respect a valid pattern (example 'app-example-log-warn-2013.09.20').

If some of your indices has a non pattern based name (like '.kibana' or 'logstash'), this tool ignores it. If you enable the log verbose option (config file), you will see that indices being (example: '2019/11/11 00:01:02 Index Date not parsed for index (258): .kibana_1').

#### DELETE SINGLE INDEX
[!WARNING] This option ignores Dry Run since the user have to specify the index to be deleted.

Usage: elasticman --delete <index_name>

This option deletes the index with the defined name. It is the same as running the curl to the API: 'curl -X DELETE "host:9200/<index_name>"'
