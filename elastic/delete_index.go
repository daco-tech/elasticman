package elastic

import (
	"elasticman/general"
	"io/ioutil"
	"log"
	"net/http"
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
