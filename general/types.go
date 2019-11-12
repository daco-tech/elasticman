package general

import "time"

type Config struct {
	Elasticsearch struct {
		Host                    string `json:"host"`
		RequiredStatus          string `json:"required_status"`
		MaxNumberOfPendingTasks int    `json:"max_number_of_pending_tasks"`
	} `json:"elasticsearch"`
	Log struct {
		Verbose bool `json:"verbose"`
	} `json:"log"`
	Parser struct {
		DateFormat         string   `json:"date_format"`
		DateIndexLastChars int      `json:"date_index_last_chars"`
		Loglevels          []string `json:"loglevels"`
	} `json:"parser"`
	Actions struct {
		Delete []struct {
			Loglevel string `json:"loglevel"`
			KeepDays int    `json:"keep-days"`
		} `json:"delete"`
	} `json:"actions"`
}

type Cluster struct {
	ClusterName                 string  `json:"cluster_name"`
	Status                      string  `json:"status"`
	TimedOut                    bool    `json:"timed_out"`
	NumberOfNodes               int     `json:"number_of_nodes"`
	NumberOfDataNodes           int     `json:"number_of_data_nodes"`
	ActivePrimaryShards         int     `json:"active_primary_shards"`
	ActiveShards                int     `json:"active_shards"`
	RelocatingShards            int     `json:"relocating_shards"`
	InitializingShards          int     `json:"initializing_shards"`
	UnassignedShards            int     `json:"unassigned_shards"`
	DelayedUnassignedShards     int     `json:"delayed_unassigned_shards"`
	NumberOfPendingTasks        int     `json:"number_of_pending_tasks"`
	NumberOfInFlightFetch       int     `json:"number_of_in_flight_fetch"`
	TaskMaxWaitingInQueueMillis int     `json:"task_max_waiting_in_queue_millis"`
	ActiveShardsPercentAsNumber float64 `json:"active_shards_percent_as_number"`
}

type Index struct {
	Health         string `json:"health"`
	Status         string `json:"status"`
	Name           string `json:"index"`
	UUID           string `json:"uuid"`
	Pri            string `json:"pri"`
	Rep            string `json:"rep"`
	DocsCount      string `json:"docs.count"`
	DocsDeleted    string `json:"docs.deleted"`
	StoreSize      string `json:"store.size"`
	PriStoreSize   string `json:"pri.store.size"`
	ParsedLogLevel string
	ParsedDate     time.Time
	ParseErrors    bool
}
