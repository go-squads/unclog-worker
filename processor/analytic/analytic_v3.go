package analytic

import (
	"github.com/Shopify/sarama"
	"github.com/go-squads/unclog-worker/filter"
	"github.com/go-squads/unclog-worker/models"
	"github.com/go-squads/unclog-worker/processor"
	"github.com/go-squads/unclog-worker/processor/alertic"
	"github.com/jasonlvhit/gocron"
	"github.com/prometheus/common/log"
	"github.com/spf13/viper"
	"strconv"
	"strings"
)

type (
	AnalyticV3Processor struct {
		repository      LogLevelMetricRepository
		analyticHandler processor.StreamHandler
	}
)

var (
	timberWolvesV3 []models.TimberWolf
	logs           map[string]map[string]int
)

func NewAnalyticV3Processor(repository LogLevelMetricRepository) (p *AnalyticV3Processor) {
	return &AnalyticV3Processor{
		repository: repository,
		analyticHandler: func(s *filter.StreamProcessorService, m *sarama.ConsumerMessage) {
			timberWolf, err := filter.ConvertKafkaMessageToTimberWolf(m)
			if err != nil {
				log.Warn(err.Error())
				return
			}

			if logs[timberWolf.ApplicationName] == nil {
				logs[timberWolf.ApplicationName] = map[string]int{}
				logs[timberWolf.ApplicationName][timberWolf.LogLevel]++
			} else {
				logs[timberWolf.ApplicationName][timberWolf.LogLevel]++
			}

			log.Info(logs)
		},
	}
}

func (p *AnalyticV3Processor) Start() {
	logs = map[string]map[string]int{}
	log.Infof("Starting Alerting System...")

	config := getAlertConfigs()

	for idx, _ := range config {
		singleConf := config[idx]

		logLevel := singleConf[0]
		appName := singleConf[1]
		interval, _ := strconv.ParseUint(singleConf[2], 10, 64)
		threshold, _ := strconv.Atoi(singleConf[3])

		gocron.Every(interval).Seconds().Do(task, appName, logLevel, threshold)

	}

	gocron.Start()
}

func getAlertConfigs() (config [][]string) {
	stringConfig := viper.GetString("ALERTS")
	splittedConfig := strings.Split(stringConfig, "|")

	for idx, _ := range splittedConfig {
		config = append(config, strings.Split(splittedConfig[idx], ","))
	}
	return
}

func (p *AnalyticV3Processor) Stop() {
	c.Stop()
}

func task(appName, logLevel string, threshold int) {
	if logs[appName][logLevel] > threshold {
		alertic.SendAlert(appName, logLevel, logs[appName][logLevel], "email")
		logs[appName][logLevel] = 0
	}
}

func (p *AnalyticV3Processor) GetHandler() processor.StreamHandler {
	return p.analyticHandler
}
