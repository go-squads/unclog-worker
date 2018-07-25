package filter

import (
	"fmt"

	"github.com/BaritoLog/go-boilerplate/errkit"
)

const (
	InvalidContextError = errkit.Error("Invalid Context Error")
	MissingContextError = errkit.Error("Missing Context Error")
)

type TimberWolf map[string]interface{}

type TimberWolfContext struct {
	KafkaTopic             string `json:"kafka_topic"`
	KafkaPartition         int32  `json:"kafka_partition"`
	KafkaReplicationFactor int16  `json:"kafka_replication_factor"`
	AppMaxTPS              int    `json:"app_max_tps"`
}

func NewTimberWolf() TimberWolf {
	timberWolf := make(map[string]interface{})
	return timberWolf
}

func (t TimberWolf) SetTimestamp(timestamp string) {
	t["@timestamp"] = timestamp
}

func (t TimberWolf) Timestamp() (s string) {
	s, _ = t["@timestamp"].(string)
	return
}

func (t TimberWolf) Context() (ctx *TimberWolfContext) {
	ctx, _ = t["_ctx"].(*TimberWolfContext)
	return
}

func (t TimberWolf) SetContext(ctx *TimberWolfContext) {
	t["_ctx"] = ctx
}

func (t TimberWolf) InitContext() (err error) {
	ctxMap, ok := t["_ctx"].(map[string]interface{})
	if !ok {
		err = MissingContextError
		return
	}

	ctx, err := mapToContext(ctxMap)
	if err != nil {
		err = errkit.Concat(InvalidContextError, err)
		return
	}

	t.SetContext(ctx)
	return
}

func mapToContext(m map[string]interface{}) (ctx *TimberWolfContext, err error) {
	kafkaTopic, ok := m["kafka_topic"].(string)
	if !ok {
		err = fmt.Errorf("kafka_topic is missing")
		return
	}

	kafkaPartition, ok := m["kafka_partition"].(float64)
	if !ok {
		err = fmt.Errorf("kafka_partition is missing")
		return
	}

	kafkaReplicationFactor, ok := m["kafka_replication_factor"].(float64)
	if !ok {
		err = fmt.Errorf("kafka_replication_factor is missing")
		return
	}

	appMaxTPS, ok := m["app_max_tps"].(float64)
	if !ok {
		err = fmt.Errorf("app_max_tps is missing")
		return
	}

	ctx = &TimberWolfContext{
		KafkaTopic:             kafkaTopic,
		KafkaPartition:         int32(kafkaPartition),
		KafkaReplicationFactor: int16(kafkaReplicationFactor),
		AppMaxTPS:              int(appMaxTPS),
	}

	return
}
