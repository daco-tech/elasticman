package elastic

import (
	"bytes"
	"elasticman/general"
	"io/ioutil"
	"log"
	"net/http"
)

// Reindex function is used to rename an index by providing the elastic endpoint, the original index name and the new index name.
// Set verbose true if you want more output details.
func Reindex(endpoint string, originalIndexName string, newIndexName string, verbose bool) bool {

	if originalIndexName != "" && newIndexName != "" && endpoint != "" {
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
		req, err := http.NewRequest("POST", endpoint+"/_reindex", bytes.NewBuffer(jsonStr))
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
		}
		if verbose {
			log.Println("response Status : ", resp.Status)
			log.Println("response Headers : ", resp.Header)
			log.Println("response Body : ", string(respBody))
		}
	}
	return false
}
