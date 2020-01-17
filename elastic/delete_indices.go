package elastic

import (
	"elasticman/general"
	types "elasticman/general"
	"elasticman/singleton"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

// DeleteIndex function is used to delete indexes by providing the elastic endpoint and the index name.
// Set verbose true if you want more output details.
func DeleteIndex(index string) bool {

	if index != "" && singleton.GetConfig().Elasticsearch.Host != "" {
		client := &http.Client{}
		req, err := http.NewRequest("DELETE", singleton.GetConfig().Elasticsearch.Host+"/"+index, nil)
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
		}
		if singleton.GetConfig().Log.Verbose {
			log.Println("response Status : ", resp.Status)
			log.Println("response Headers : ", resp.Header)
			log.Println("response Body : ", string(respBody))
		}
	}
	return false
}

// DeleteByDays function is used to delete indexes by days. This function also provides a dry run option.
// It uses a list of indexes provided by the GetParsedIndices function.
// Set verbose true if you want more output details.
func DeleteByDays(parsedIndices []types.Index, days int, logtype string, loglevel string) int {
	var deletedIndices int
	var possibleDeletions int
	for _, index := range parsedIndices {
		if !index.ParseErrors {
			var delete int
			if loglevel != "" && loglevel == index.ParsedLogLevel {
				delete++
			}
			if (logtype != "" && logtype == index.ParsedLogType) || logtype == "" {
				delete++
			}
			if index.ExistenceInDays >= days {
				delete++
			}
			if delete == 3 {
				if !singleton.GetConfig().Actions.Delete.DryRun {
					DeleteIndex(index.Name)
					deletedIndices++
					log.Println("Index with name '" + index.Name + "' has been deleted since it has " + strconv.Itoa(index.ExistenceInDays) + " days and logtype/loglevel '" + index.ParsedLogType + "'/'" + index.ParsedLogLevel + "'.")
				} else {
					possibleDeletions++
					log.Println("*** DRY RUN *** - Index with name '" + index.Name + "' could been deleted since it has " + strconv.Itoa(index.ExistenceInDays) + " days and logtype/loglevel '" + index.ParsedLogType + "'/'" + index.ParsedLogLevel + "'.")
				}
			}
		}
	}

	var text string
	if singleton.GetConfig().Actions.Delete.DryRun {
		text = text + " Possible deletions: " + strconv.Itoa(possibleDeletions)
	} else {
		if deletedIndices > 0 {
			text = "Deleted Indexes (" + loglevel + "): " + strconv.Itoa(deletedIndices)
		} else {
			if singleton.GetConfig().Log.Verbose {
				text = "Nothing deleted (" + loglevel + ")!"
			}
		}
	}
	if text != "" {
		log.Println(text)
	}
	return deletedIndices
}
