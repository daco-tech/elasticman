[![Codacy Badge](https://api.codacy.com/project/badge/Grade/5f66fc949da24d148b2b7bb274d347e4)](https://www.codacy.com/manual/daco-tech/elasticman?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=daco-tech/elasticman&amp;utm_campaign=Badge_Grade) ![](https://github.com/daco-tech/elasticman/workflows/ElasticMan-Build/badge.svg)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fdaco-tech%2Felasticman.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Fdaco-tech%2Felasticman?ref=badge_shield)

# ElasticMan - Elastic Maintenance Tool

This is a simple tool to maintain the ElasticSearch Log indices.

The main objective of this project is to create a tool with enough flexibility to clean, consolidate and organize ElasticSearch indices easily.

The intent of this tool it's not to compete with ElasticSearch Curator, that is a more complete and maintained tool, but it is a tool that you can keep on your laptop or even in a docker container with very low hardware resources usage to allow you to maintain and quickly respond to the day to day challenges.

## Getting started

### From Source

*   Make sure you have [dep](https://github.com/golang/dep) installed
*   Clone this repo `git clone https://github.com/daco-tech/elasticman.git`
*   Create a configuration file with name: config.json at main.go file level with the text in the config section
*   Run `make` to download dependencies and run the application

## Config

By default, elasticman searches for the configuration file in a file called config.json in the .elasticsearch directory inside your home directory.

You can override this path by passing -c or --config option (i.e.: elasticman --config ./clusterA.json ...).

~/.elasticman/config.json file content example:

``` json
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
            "critical",
            "severe",
            "fatal"
        ],
        "logtypes": [
            "log",
            "evt",
            "orders",
            "metric"
        ],
        "ignorelist": [
            "trash",
            "^.kibana[a-zA-Z0-9_.]*$",
            "kibana_backup",
            "logstash",
            ".*(-all)$"
        ]
    },
    "actions": {
        "delete": {
            "enabled": true,
            "dry_run": false,
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
                    "loglevel": "severe",
                    "keep-days": 60
                },
                {
                    "loglevel": "fatal",
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
Usage: elasticman delete

This command deletes indices according with the configuration defined in the actions -> delete section. 

This tool can be tested without delete indices if the  "dry_run" option is set to true in the configuration file.

Once you test the configuration file, you are able to schedule in a cron job the automatic execution of this cleanup tool to keep your elastic search indices clean (ideal for logs).

Be aware that the indices date should be in the end of their name and should respect a valid pattern (example 'app-example-log-warn-2013.09.20').

If some of your indices has a non pattern based name (like '.kibana' or 'logstash'), this tool ignores it. If you enable the log verbose option (config file), you will see that indices being (example: '2019/11/11 00:01:02 Index Date not parsed for index (258): .kibana_1').

#### DELETE SINGLE INDEX
**WARNING** - This option ignores Dry Run since the user have to specify the index to be deleted.

Usage: elasticman delete --index <index_name>

This option deletes the index with the defined name. It is the same as running the curl to the API: 'curl -X DELETE "host:9200/<index_name>"'


### REINDEX

Usage: elasticman reindex \<indexname\> \<new_or_existing_index\>

Reindexes index data into a new index name or to an existing index.

This option does not delete the original index, if you want to rename an existing index you should reindex to the new name and then delete the original index after verify that everything was reindexed to the new one.

This option is also useful if you changed the matching index template and you need to reindex the index to meet the new one.

### CONSOLIDATE

Usage: elasticman consolidate

Reindexes index data into a new configured consolidation destination index.
This option deletes the original index after reindex all the documents to the destination index.

You need to provide the necessary consolidation configuration at the config.json file.

Example:

```
{
(...)
"actions": {
    (...)
    "consolidate": {
            "enabled": true,
            "dry_run": false,
            "delete_source_index": true,
            "todo": [
                {
                    "loglevel": "info",
                    "keep-days": 1,
                    "logtype": "evt",
                    "add-suffix": "all",
                    "add-current-month": true,
                    "add-current-year": true,
                    "remove-log-level": true,
                    "remove-date": true,
                    "remove-last-no-of-chars": 0
                }
            ]
        }
    }
}

```
