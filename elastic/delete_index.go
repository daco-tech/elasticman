package elastic

import (
	"elasticman/general"
	types "elasticman/general"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

func DeleteIndex(endpoint string, index string, verbose bool) bool {

	if index != "" && endpoint != "" {
		client := &http.Client{}
		req, err := http.NewRequest("DELETE", endpoint+"/"+index, nil)
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
		if verbose {
			log.Println("response Status : ", resp.Status)
			log.Println("response Headers : ", resp.Header)
			log.Println("response Body : ", string(respBody))
		}
	}
	return false
}

func DeleteByDays(endpoint string, dryrun bool, parsedIndices []types.Index, days int, logtype string, loglevel string, verbose bool) bool {
	for _, index := range parsedIndices {
		if !index.ParseErrors {
			var delete int = 0
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
				if !dryrun {
					DeleteIndex(endpoint, index.Name, verbose)
					log.Println("Index with name '" + index.Name + "' has been deleted since it has " + strconv.Itoa(index.ExistenceInDays) + " days and logtype/loglevel '" + index.ParsedLogType + "'/'" + index.ParsedLogLevel + "'.")
				} else {
					log.Println("*** DRY RUN *** - Index with name '" + index.Name + "' could been deleted since it has " + strconv.Itoa(index.ExistenceInDays) + " days and logtype/loglevel '" + index.ParsedLogType + "'/'" + index.ParsedLogLevel + "'.")
				}
			}
		}

	}
	return false
}
