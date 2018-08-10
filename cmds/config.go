package cmds

import (
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

const (
	EnvKafkaBrokers      = "BARITO_KAFKA_BROKERS"
	EnvKafkaGroupID      = "BARITO_KAFKA_GROUP_ID"
	EnvKafkaTopicSuffix  = "BARITO_KAFKA_TOPIC_SUFFIX"
	EnvConsulUrl         = "BARITO_CONSUL_URL"
	EnvConsulKafkaName   = "BARITO_CONSUL_KAFKA_NAME"
	EnvNewTopicEventName = "BARITO_NEW_TOPIC_EVENT"
)

var (
	DefaultConsulKafkaName   = "kafka"
	DefaultKafkaBrokers      = []string{"localhost:9092"}
	DefaultKafkaTopicSuffix  = "_logs"
	DefaultKafkaGroupID      = "barito-group"
	DefaultNewTopicEventName = "new_topic_events"
)

func configKafkaBrokers() (brokers []string) {
	consulUrl := configConsulUrl()
	name := configConsulKafkaName()
	brokers, err := consulKafkaBroker(consulUrl, name)
	if err != nil {
		brokers = sliceEnvOrDefault(EnvKafkaBrokers, ",", DefaultKafkaBrokers)
		return
	}

	logConfig("consul", EnvKafkaBrokers, brokers)
	return
}

func configKafkaGroupId() (s string) {
	return stringEnvOrDefault(EnvKafkaGroupID, DefaultKafkaGroupID)
}

func configConsulKafkaName() (s string) {
	return stringEnvOrDefault(EnvConsulKafkaName, DefaultConsulKafkaName)
}

func configConsulUrl() (s string) {
	return os.Getenv(EnvConsulUrl)
}

func configKafkaTopicSuffix() string {
	return stringEnvOrDefault(EnvKafkaTopicSuffix, DefaultKafkaTopicSuffix)
}

func configNewTopicEvent() string {
	return stringEnvOrDefault(EnvNewTopicEventName, DefaultNewTopicEventName)

}

func stringEnvOrDefault(key, defaultValue string) string {
	s := os.Getenv(key)
	if len(s) > 0 {
		logConfig("env", key, s)
		return s
	}

	logConfig("default", key, defaultValue)
	return defaultValue
}

func sliceEnvOrDefault(key, separator string, defaultSlice []string) []string {
	s := os.Getenv(key)

	if len(s) > 0 {
		slice := strings.Split(s, separator)
		logConfig("env", key, slice)
		return slice
	}

	logConfig("default", key, defaultSlice)
	return defaultSlice
}

func logConfig(source, key string, val interface{}) {
	log.WithField("config", source).Infof("%s = %v", key, val)
}
