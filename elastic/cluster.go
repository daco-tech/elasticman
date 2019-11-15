package elastic

import (
	"elasticman/general"
	types "elasticman/general"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

// GetClusterStatus function returns the cluster status details.
// Set verbose true if you want more output details.
func GetClusterStatus(endpoint string, verbose bool) (types.Cluster, string) {
	clusterStatus := types.Cluster{}
	if endpoint != "" {
		client := &http.Client{}
		req, err := http.NewRequest("GET", endpoint+"/_cluster/health", nil)
		if err != nil {
			log.Fatalln(err)
			return clusterStatus, "not connected"
		}

		// Fetch Request
		resp, err := client.Do(req)
		if err != nil {
			log.Fatalln(err)
			return clusterStatus, err.Error()
		}
		defer resp.Body.Close()

		// Read Response Body
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
			return clusterStatus, err.Error()
		}

		// Display Results
		if general.HasPrefix(resp.Status, "200") {

			jsonErr := json.Unmarshal(respBody, &clusterStatus)
			if jsonErr != nil {
				return clusterStatus, "JSON parse error."
			}
			return clusterStatus, ""
		}
		if verbose {
			log.Println("response Status : ", resp.Status)
			log.Println("response Headers : ", resp.Header)
			log.Println("response Body : ", string(respBody))
		}

	} else {
		return clusterStatus, "not connected"
	}
	return clusterStatus, "not configured"
}
