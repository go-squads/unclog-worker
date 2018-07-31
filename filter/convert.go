package filter

import (
	"encoding/json"
	"time"

	"github.com/BaritoLog/go-boilerplate/errkit"
	"github.com/Shopify/sarama"
)

const (
	JsonParseError = errkit.Error("JSON Parse Error")
)

func ConvertBytesToTimberWolf(data []byte) (timberWolf TimberWolf, err error) {
	err = json.Unmarshal(data, &timberWolf)
	if err != nil {
		err = errkit.Concat(JsonParseError, err)
		return
	}

	err = timberWolf.InitContext()
	if err != nil {
		return
	}

	if timberWolf.Timestamp() == "" {
		timberWolf.SetTimestamp(time.Now().UTC().Format(time.RFC3339))
	}

	if timberWolf.LogLevel() == "" {
		timberWolf.SetLogLevel("UNLISTED")
	}

	cleanTimberWolf(timberWolf)

	return
}

func cleanTimberWolf(timberWolf TimberWolf) {
	attributes := []string{"_ctx", "log_level", "@timestamp"}

	doc := make(map[string]interface{})
	for k, v := range timberWolf {
		doc[k] = v

		if !stringInSlice(k, attributes) {
			delete(timberWolf, k)
		}
	}
}

// NewTimberWolfFromKafkaMessage create timberWolf instance from kafka message
func ConvertKafkaMessageToTimberWolf(message *sarama.ConsumerMessage) (timberWolf TimberWolf, err error) {
	return ConvertBytesToTimberWolf(message.Value)
}

// ConvertToKafkaMessage will convert timberWolf to sarama producer message for kafka
func ConvertTimberWolfToKafkaMessage(timberWolf TimberWolf, topic string) *sarama.ProducerMessage {
	b, _ := json.Marshal(timberWolf)

	return &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(b),
	}
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
