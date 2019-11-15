package general

import (
	"encoding/json"
	"os"
)

// LoadConfiguration Loads the configuration file to the Config type.
// Set verbose true if you want more output details.
func LoadConfiguration(filename string) (Config, error) {
	var config Config
	configFile, err := os.Open(filename)
	defer configFile.Close()
	if err != nil {
		return config, err
	}
	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&config)
	return config, err
}
