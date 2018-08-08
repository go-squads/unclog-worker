package filter

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/BaritoLog/go-boilerplate/errkit"
	"github.com/Shopify/sarama"
	"github.com/go-squads/unclog-worker/models"
)

const (
	JsonParseError = errkit.Error("JSON Parse Error")
)

func ConvertBytesToTimberWolf(data []byte) (timberWolf models.TimberWolf, err error) {
	mappedData := make(map[string]interface{})

	err = json.Unmarshal(data, &mappedData)
	if err != nil {
		err = errkit.Concat(JsonParseError, err)
		return
	}

	timberWolf = models.TimberWolf{}

	if mappedData["@timestamp"] == nil || mappedData["@timestamp"] == "" {
		timberWolf.Timestamp = time.Now().UTC().Format("2006-01-02 15:04:05")
	} else {
		timberWolf.Timestamp = mappedData["@timestamp"].(string)
	}

	if mappedData["log_level"] == nil || mappedData["log_level"] == "" {
		timberWolf.LogLevel = "UNLISTED"
	} else {
		timberWolf.LogLevel = strings.ToUpper(mappedData["log_level"].(string))
	}

	if mappedData["app_name"] == nil || mappedData["app_name"] == "" {
		timberWolf.ApplicationName = "UNKNOWN APP"
	} else {
		timberWolf.ApplicationName = mappedData["app_name"].(string)
	}

	if mappedData["node_id"] == nil || mappedData["node_id"] == "" {
		timberWolf.NodeId = "UNKNOWN NODE"
	} else {
		timberWolf.NodeId = mappedData["node_id"].(string)
	}

	return
}

// NewTimberWolfFromKafkaMessage create timberWolf instance from kafka message
func ConvertKafkaMessageToTimberWolf(message *sarama.ConsumerMessage) (timberWolf models.TimberWolf, err error) {
	return ConvertBytesToTimberWolf(message.Value)
}

// ConvertToKafkaMessage will convert timberWolf to sarama producer message for kafka
func ConvertTimberWolfToKafkaMessage(timberWolf models.TimberWolf, topic string) *sarama.ProducerMessage {
	b, _ := json.Marshal(timberWolf)

	return &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(b),
	}
}
