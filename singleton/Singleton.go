package singleton

import (
	"elasticman/general"
	"sync"
)

//singleton variables
type singleton map[string]interface{}

var (
	// ensures that our type only gets initialized exactly once.
	once sync.Once

	instance singleton
)

//init once
func init() {
	once.Do(func() {

		instance = make(singleton)
	})
}

//Set key-value
func Set(key string, value interface{}) {
	instance[key] = value
}

//Get return value
func Get(key string) interface{} {
	return instance[key]
}

func LoadConfig(configFile string) {

}

func SetConfig(config general.Config) {
	instance["config"] = config
}
func GetConfig() general.Config {
	return instance["config"].(general.Config)
}
