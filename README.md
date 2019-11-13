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