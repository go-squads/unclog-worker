package analytic

import (
	"encoding/json"
	"github.com/Shopify/sarama"
	"github.com/go-squads/unclog-worker/filter"
	"github.com/go-squads/unclog-worker/models"
	"github.com/go-squads/unclog-worker/processor"
	"github.com/go-squads/unclog-worker/processor/alertic"
	"github.com/jasonlvhit/gocron"
	"github.com/prometheus/common/log"
)

type (
	AnalyticV3Processor struct {
		repository      LogLevelMetricRepository
		analyticHandler processor.StreamHandler
	}

	AlertConfig struct {
		Id       int
		AppName  string
		LogLevel string
		Duration int
		Limit    int
		Callback string
	}
)

var (
	timberWolvesV3 []models.TimberWolf
	logs           map[string]map[string]int
	cronJobV3      *gocron.Scheduler
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

	cronJobV3 = gocron.NewScheduler()
	config, err := p.repository.GetAlertConfig()

	if err != nil {
		log.Error(err)
	}

	for idx, _ := range config {
		a := config[idx]

		logLevel := a.LogLevel
		appName := a.AppName
		duration := uint64(a.Duration)
		limit := a.Limit

		var cb map[string]interface{}
		err = json.Unmarshal([]byte(a.Callback), &cb)
		if err != nil {
			log.Fatal(err)
		}

		callbackType := cb["type"]
		callbackReceiver := convertInterfaceToString(cb["receivers"].([]interface{}))

		cronJobV3.Every(duration).Seconds().Do(task, appName, logLevel, callbackType, callbackReceiver, limit)
	}

	cronJobV3.Start()
}

func convertInterfaceToString(arrayOfInterface []interface{}) (receivers []string) {
	for idx, _ := range arrayOfInterface {
		receivers = append(receivers, arrayOfInterface[idx].(string))
	}

	return
}

func (p *AnalyticV3Processor) Stop() {
	cronJobV3.Clear()
}

func task(appName, logLevel, callbackType string, callbackReceivers []string, threshold int) {
	if logs[appName][logLevel] > threshold {
		alertic.SendAlert(appName, logLevel, logs[appName][logLevel], callbackType, callbackReceivers)
		logs[appName][logLevel] = 0
	}
}

func (p *AnalyticV3Processor) GetHandler() processor.StreamHandler {
	return p.analyticHandler
}
