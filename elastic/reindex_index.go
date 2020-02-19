package elastic

import (
	"bytes"
	"elasticman/general"
	types "elasticman/general"
	"elasticman/singleton"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Reindex function is used to rename an index by providing the elastic endpoint, the original index name and the new index name.
// Set verbose true if you want more output details.
func Reindex(originalIndexName string, newIndexName string) bool {

	if originalIndexName != "" && newIndexName != "" && singleton.GetConfig().Elasticsearch.Host != "" {
		client := &http.Client{}
		var jsonStr = []byte(`
		{
			"source": {
				"index": "` + originalIndexName + `"
			},
			"dest": {
				"index": "` + newIndexName + `"
			}
		}`)
		req, err := http.NewRequest("POST", singleton.GetConfig().Elasticsearch.Host+"/_reindex", bytes.NewBuffer(jsonStr))
		req.Header.Set("Content-Type", "application/json")
		if err != nil {
			log.Fatalln(err)
			return false
		}

		// Fetch Request
		resp, err := client.Do(req)
		if err != nil {
			log.Fatalln(err)
			return false
		}
		defer resp.Body.Close()

		// Read Response Body
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
			return false
		}

		// Display Results
		if general.HasPrefix(resp.Status, "200") {
			return true
		} else {
			log.Println("Error reindexing index. Request: " + string(jsonStr) + " --- REST Request Status: " + resp.Status)
		}
		if singleton.GetVerbose() {
			log.Println("response Status : ", resp.Status)
			log.Println("response Headers : ", resp.Header)
			log.Println("response Body : ", string(respBody))
		}
	}
	return false
}

func Rename(originalIndexName string, newIndexName string) bool {
	result := Reindex(originalIndexName, newIndexName)
	if result {
		result = DeleteIndex(originalIndexName)
	} else {
		log.Println("Error Re-Indexing. Source index with name '" + originalIndexName + "' was not deleted! Rename Failed!")
	}
	return result
}

// ConsolidateByDays function is used to merge indexes by days. This function also provides a dry run option.
// It uses a list of indexes provided by the GetParsedIndices function.
// Set verbose true if you want more output details.
func ConsolidateByDays(parsedIndices []types.Index, days int, logtype string, loglevel string, suffix string, dryrun bool, delete bool, addyear bool, removeloglevel bool, removeDate bool) int {
	var consolidatedIndices int
	var possibleConsolidations int
	for _, index := range parsedIndices {
		if !index.ParseErrors {

			if logtype == index.ParsedLogType && loglevel == index.ParsedLogLevel {
				var consolidate bool
				if index.ExistenceInDays >= days {
					consolidate = true
				}
				newIndexName := index.Name
				if removeDate {
					newIndexName = index.Name[:len(index.Name)-(singleton.GetConfig().Parser.DateIndexLastChars+1)]
				}
				if removeloglevel {
					newIndexName = strings.Replace(newIndexName, "-"+loglevel, "", 1)

				}
				if suffix != "" {
					newIndexName += "-" + suffix
				}
				if addyear {
					currentTime := time.Now()
					newIndexName += "-" + strconv.Itoa(currentTime.Year())
				}
				if newIndexName == index.Name {
					log.Fatal("Error! Configured destination ('" + newIndexName + "') consolidation index is the same of the source index ('" + index.Name + "')!")
				}

				if consolidate {
					if !singleton.GetConfig().Actions.Consolidate.DryRun {
						var result bool
						if delete {
							result = Rename(index.Name, newIndexName)
						} else {
							result = Reindex(index.Name, newIndexName)
						}
						if result {
							consolidatedIndices++
							log.Println("Index with name '" + index.Name + "' has been merged since it has " + strconv.Itoa(index.ExistenceInDays) + " days and logtype/loglevel '" + index.ParsedLogType + "'/'" + index.ParsedLogLevel + "' onto the index '" + newIndexName + "'.")
						} else {
							consolidatedIndices--
							log.Println("Error merging index with name: '" + index.Name + "'.")
						}

					} else {
						possibleConsolidations++
						log.Println("*** DRY RUN *** - Index with name '" + index.Name + "' could been merged since it has " + strconv.Itoa(index.ExistenceInDays) + " days and logtype/loglevel '" + index.ParsedLogType + "'/'" + index.ParsedLogLevel + "' onto the index '" + newIndexName + "'.")
					}
				}
			}
		}
	}

	// Final Resume
	var text string
	if singleton.GetConfig().Actions.Consolidate.DryRun {
		text = text + " Possible merges: " + strconv.Itoa(possibleConsolidations)
	} else {
		if consolidatedIndices > 0 {
			text = "MergedIndexes (" + loglevel + "): " + strconv.Itoa(consolidatedIndices)
		} else {
			if singleton.GetVerbose() {
				text = "Nothing Merged (" + loglevel + ")!"
			}
		}
	}
	if text != "" {
		log.Println(text)
	}
	return consolidatedIndices
}
