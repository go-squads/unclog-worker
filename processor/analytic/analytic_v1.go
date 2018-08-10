package analytic

import (
	"github.com/Shopify/sarama"
	"github.com/go-squads/unclog-worker/filter"
	"github.com/go-squads/unclog-worker/models"
	"github.com/go-squads/unclog-worker/processor"
	"github.com/prometheus/common/log"
	"github.com/robfig/cron"
	"github.com/spf13/viper"
)

type (
	AnalyticV1Processor struct {
		repository      LogLevelMetricRepository
		analyticHandler processor.StreamHandler
	}
)

var timberWolvesV1 []models.TimberWolf

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
	interval := viper.GetString("WINDOWING_INTERVAL_IN_SECOND_V1") + "s"

	c := cron.New()
	c.AddFunc("@every "+interval, func() {
		p.saveToDatabase()
		log.Info(timberWolvesV1)
		timberWolvesV1 = []models.TimberWolf{}
	})
	go c.Start()
}

func (p *AnalyticV1Processor) Stop() {
	c.Stop()
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
