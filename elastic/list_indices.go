package elastic

import (
	"elasticman/general"
	types "elasticman/general"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/metakeule/fmtdate"
)

// This function returns indices list from ElasticSearch.
// Set verbose true if you want more output details.
func GetIndices(endpoint string, verbose bool) ([]types.Index, string) {
	indices := []types.Index{}

	if endpoint != "" {
		client := &http.Client{}
		req, err := http.NewRequest("GET", endpoint+"/_cat/indices?format=json&bytes=b", nil)
		if err != nil {
			log.Fatalln(err)
			return indices, "not connected"
		}

		// Fetch Request
		resp, err := client.Do(req)
		if err != nil {
			log.Fatalln(err)
			return indices, err.Error()
		}
		defer resp.Body.Close()

		// Read Response Body
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
			return indices, err.Error()
		}

		// Display Results
		if general.HasPrefix(resp.Status, "200") {

			jsonErr := json.Unmarshal(respBody, &indices)
			if jsonErr != nil {
				return indices, "JSON parse error."
			}
			log.Println("Returning indices...")
			return indices, ""
		}
		if verbose {
			log.Println("response Status : ", resp.Status)
			log.Println("response Headers : ", resp.Header)
			log.Println("response Body : ", string(respBody))
		}

	} else {
		return indices, "not connected"
	}
	return indices, "not configured"
}

// This function returns indices list with parsed fields filled. Like the existence days of the index, loglevel, logtype and index date.
// Set verbose true if you want more output details.
func GetParsedIndices(endpoint string, verbose bool, dateformat string, dateLastNoOfChars int, loglevels []string, logtypes []string) (parsedIndices []types.Index, err string) {
	var indices, getErr = GetIndices(endpoint, verbose)
	parsedIndices = []types.Index{}
	if getErr == "" {
		for i, index := range indices {
			var indexMod = indices[i]
			//Parse Date
			if len(index.Name) > dateLastNoOfChars {
				var data string = string(index.Name[len(index.Name)-dateLastNoOfChars:])
				t, parseErr := fmtdate.Parse(dateformat, data)
				if parseErr == nil {
					indexMod.ParsedDate = t
					//Calculate Days of Existence
					indexMod.ExistenceInDays = int(time.Now().Sub(indexMod.ParsedDate).Hours() / 24)
				} else {
					indexMod.ParseErrors = true
				}
			}

			//Parse LogLevel
			if len(loglevels) > 0 {
				for _, loglevel := range loglevels {
					if strings.Contains(index.Name, loglevel) {
						indexMod.ParsedLogLevel = loglevel
					}
				}
			}

			//Parse LogTypes
			if len(logtypes) > 0 {
				for _, logtype := range logtypes {
					if strings.Contains(index.Name, logtype) {
						indexMod.ParsedLogType = logtype
					}
				}
			}

			//Advice Parse Issues
			if indexMod.ParsedDate.IsZero() {
				indexMod.ParseErrors = true
				if verbose {
					log.Println("Index Date not parsed for index (" + strconv.Itoa(i) + "): " + indexMod.Name)
				}
			}

			if indexMod.ParsedLogLevel == "" && len(loglevels) > 0 {
				indexMod.ParseErrors = true
				if verbose {
					log.Println("LogLevel not parsed for index (" + strconv.Itoa(i) + "): " + indexMod.Name)
				}
			}

			if indexMod.ParsedLogType == "" && len(logtypes) > 0 {
				indexMod.ParseErrors = true
				if verbose {
					log.Println("LogType not parsed for index (" + strconv.Itoa(i) + "): " + indexMod.Name)
				}
			}

			parsedIndices = append(parsedIndices, indexMod)
		}
	}
	return parsedIndices, getErr
}
