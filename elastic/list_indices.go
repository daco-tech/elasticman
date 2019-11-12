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

	"github.com/metakeule/fmtdate"
)

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

func GetParsedIndices(endpoint string, verbose bool, dateformat string, date_last_no_of_chars int, loglevels []string) (parsed_indices []types.Index, err string) {
	var indices, getErr = GetIndices(endpoint, verbose)
	parsed_indices = []types.Index{}
	log.Println("GetParsedIndices")
	if getErr == "" {
		for i, index := range indices {
			var indexMod = indices[i]
			//PARSE DATE
			if len(index.Name) > date_last_no_of_chars {
				var data string = string(index.Name[len(index.Name)-date_last_no_of_chars:])
				t, parseErr := fmtdate.Parse(dateformat, data)
				if parseErr == nil {
					indexMod.ParsedDate = t
				} else {
					indexMod.ParseErrors = true
				}
			}

			//PARSE LOGLEVEL
			for _, loglevel := range loglevels {
				if strings.Contains(index.Name, loglevel) {
					indexMod.ParsedLogLevel = loglevel
				}
			}

			parsed_indices = append(parsed_indices, indexMod)
			if indexMod.ParsedDate.IsZero() && verbose {
				log.Println("Index Date not parsed for index (" + strconv.Itoa(i) + "): " + indexMod.Name)
			}

		}
	}

	return parsed_indices, getErr
}
