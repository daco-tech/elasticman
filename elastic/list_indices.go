package elastic

import (
	"elasticman/general"
	types "elasticman/general"
	"elasticman/singleton"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/metakeule/fmtdate"
)

// GetIndices function returns indices list from ElasticSearch.
// Set verbose true if you want more output details.
func GetIndices() ([]types.Index, string) {
	indices := []types.Index{}

	if singleton.GetConfig().Elasticsearch.Host != "" {
		client := &http.Client{}
		req, err := http.NewRequest("GET", singleton.GetConfig().Elasticsearch.Host+"/_cat/indices?format=json&bytes=b", nil)
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
		if singleton.GetVerbose() {
			log.Println("response Status : ", resp.Status)
			log.Println("response Headers : ", resp.Header)
			log.Println("response Body : ", string(respBody))
		}

	} else {
		return indices, "not connected"
	}
	return indices, "not configured"
}

// GetIndicesWithoutIgnored function returns indices list with parsed fields filled removing the user ignored indexes set in the config file. Like the existence days of the index, loglevel, logtype and index date.
// Set verbose true if you want more output details.
func GetIndicesWithoutIgnored() (parsedIndices []types.Index, err string) {
	var indices, getErr = GetIndices()
	cleanIndices := make([]types.Index, len(parsedIndices))
	for _, indexed := range indices {
		var ignorable bool
		for _, ignored := range singleton.GetConfig().Parser.Ignorelist {
			if ignored != "" {
				r, _ := regexp.Compile(ignored)

				if r.MatchString(indexed.Name) {

					ignorable = true
					if singleton.GetVerbose() {
						log.Println("Index name: " + indexed.Name + " matches the regex: " + ignored)
					}
					break
				}
			}
		}
		if !ignorable {
			cleanIndices = append(cleanIndices, indexed)
		}
		ignorable = false
	}
	return cleanIndices, getErr
}

// GetParsedIndices function returns indices list with parsed fields filled. Like the existence days of the index, loglevel, logtype and index date.
// Set verbose true if you want more output details.
func GetParsedIndices() (parsedIndices []types.Index, err string) {
	var indices, getErr = GetIndicesWithoutIgnored()
	log.Println("Eligible Indices to be parsed: " + strconv.Itoa(len(indices)) + ";")

	if getErr == "" {
		for i, index := range indices {
			var indexMod = indices[i]
			//Parse Date
			if len(index.Name) > singleton.GetConfig().Parser.DateIndexLastChars {
				var data string = string(index.Name[len(index.Name)-singleton.GetConfig().Parser.DateIndexLastChars:])
				t, parseErr := fmtdate.Parse(singleton.GetConfig().Parser.DateFormat, data)
				if parseErr == nil {
					indexMod.ParsedDate = t
					//Calculate Days of Existence
					indexMod.ExistenceInDays = int(time.Now().Sub(indexMod.ParsedDate).Hours() / 24)
				} else {
					indexMod.ParseErrors = true
				}
			}

			//Parse LogLevel
			if len(singleton.GetConfig().Parser.Loglevels) > 0 {
				for _, loglevel := range singleton.GetConfig().Parser.Loglevels {
					if strings.Contains(index.Name, loglevel) {
						indexMod.ParsedLogLevel = loglevel
					}
				}
			}

			//Parse LogTypes
			if len(singleton.GetConfig().Parser.Logtypes) > 0 {
				for _, logtype := range singleton.GetConfig().Parser.Logtypes {
					if strings.Contains(index.Name, logtype) {
						indexMod.ParsedLogType = logtype
					}
				}
			}

			//Advice Parse Issues
			if indexMod.ParsedDate.IsZero() {
				indexMod.ParseErrors = true
				if singleton.GetVerbose() {
					log.Println("Index Date not parsed for index (" + strconv.Itoa(i) + "): " + indexMod.Name)
				}
			}

			if indexMod.ParsedLogLevel == "" && len(singleton.GetConfig().Parser.Loglevels) > 0 {
				indexMod.ParseErrors = true
				if singleton.GetVerbose() {
					log.Println("LogLevel not parsed for index (" + strconv.Itoa(i) + "): " + indexMod.Name)
				}
			}

			if indexMod.ParsedLogType == "" && len(singleton.GetConfig().Parser.Logtypes) > 0 {
				indexMod.ParseErrors = true
				if singleton.GetVerbose() {
					log.Println("LogType not parsed for index (" + strconv.Itoa(i) + "): " + indexMod.Name)
				}
			}

			parsedIndices = append(parsedIndices, indexMod)
		}
	}
	return parsedIndices, getErr
}
