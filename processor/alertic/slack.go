package alertic

import (
	"strconv"
	"bytes"
	"net/http"
	"github.com/prometheus/common/log"
)

func buildMessage(appName, logLevel string, currentCount int) (message string) {
	return "{\"text\" :" + "\"" + appName + " had " + strconv.Itoa(currentCount) + " " + logLevel + "\"}"
}

func sendMessage(message, receiver string) {
	_, err := http.Post(receiver, "application/json", bytes.NewBuffer([]byte(message)))

	if err != nil {
		log.Info(err.Error())
	}

	log.Info("Message sent!!")
}

func SendToSlack(appName, logLevel string, currentCount int, receivers []string) {
	message := buildMessage(appName, logLevel, currentCount)

	for idx, _ := range receivers {
		receiver := receivers[idx]

		sendMessage(message, receiver)
	}
}
