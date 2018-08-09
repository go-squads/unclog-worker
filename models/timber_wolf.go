package models

type TimberWolf struct {
	Timestamp       string `json:"@timestamp"`
	ApplicationName string `json:"app_name"`
	NodeId          string `json:"node_id"`
	LogLevel        string `json:"log_level"`
	Counter         int    `json:"counter"`
}
