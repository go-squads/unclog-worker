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
	KafkaTopic string `json:"kafka_topic"`
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

func (t TimberWolf) SetLogLevel(loglevel string) {
	t["log_level"] = loglevel
}

func (t TimberWolf) LogLevel() (s string) {
	s, _ = t["log_level"].(string)
	return
}

func (t TimberWolf) LogLevelId() int {
	s := t.LogLevel()

	switch s {
	case "ERROR":
		return 1
	case "INFO":
		return 2
	case "TRACE":
		return 3
	case "DEBUG":
		return 4
	case "WARN":
		return 5
	default:
		return 0
	}

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

	ctx = &TimberWolfContext{
		KafkaTopic: kafkaTopic,
	}

	return
}
