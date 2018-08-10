package analytic

import (
	"github.com/Shopify/sarama"
	"github.com/go-squads/unclog-worker/filter"
	"github.com/go-squads/unclog-worker/models"
	"github.com/go-squads/unclog-worker/processor"
	"github.com/jasonlvhit/gocron"
	"github.com/prometheus/common/log"
	"github.com/spf13/viper"
)

type (
	AnalyticV1Processor struct {
		repository      LogLevelMetricRepository
		analyticHandler processor.StreamHandler
	}
)

var (
	timberWolvesV1 []models.TimberWolf
	cronJobV1      *gocron.Scheduler
)

func NewAnalyticV1Processor(repository LogLevelMetricRepository) (p *AnalyticV1Processor) {
	return &AnalyticV1Processor{
		repository: repository,
		analyticHandler: func(s *filter.StreamProcessorService, m *sarama.ConsumerMessage) {
			timberWolf, err := filter.ConvertKafkaMessageToTimberWolf(m)
			if err != nil {
				log.Warn(err.Error())
				return
			}

			if idx := getTimberWolfV1Index(timberWolf); idx != -1 {
				timberWolvesV1[idx].Counter++
			} else {
				timberWolvesV1 = append(timberWolvesV1, timberWolf)
			}

			log.Info(timberWolf)
		},
	}
}

func (p *AnalyticV1Processor) Start() {
	log.Infof("Starting Analytic Processor...")
	interval := viper.GetInt("WINDOWING_INTERVAL_IN_SECOND_V1")

	cronJobV1 = gocron.NewScheduler()

	cronJobV1.Every(uint64(interval)).Seconds().Do(func() {
		p.saveToDatabase()
		log.Info(timberWolvesV1)
		timberWolvesV1 = []models.TimberWolf{}
	})
	cronJobV1.Start()
}

func (p *AnalyticV1Processor) Stop() {
	cronJobV1.Clear()
}

func (p *AnalyticV1Processor) saveToDatabase() {
	for _, timberWolf := range timberWolvesV1 {
		p.repository.SaveV1(timberWolf)
	}

	log.Info("Logs aggregation stored!")
}

func (p *AnalyticV1Processor) GetHandler() processor.StreamHandler {
	return p.analyticHandler
}

func getTimberWolfV1Index(t models.TimberWolf) int {
	for idx, timberWolf := range timberWolvesV1 {
		if timberWolf.LogLevel == t.LogLevel {
			return idx
		}
	}

	return -1
}
