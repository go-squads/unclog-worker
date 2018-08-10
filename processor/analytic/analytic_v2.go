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
	AnalyticV2Processor struct {
		repository      LogLevelMetricRepository
		analyticHandler processor.StreamHandler
	}
)

var timberWolvesV2 []models.TimberWolf

func NewAnalyticV2Processor(repository LogLevelMetricRepository) (p *AnalyticV2Processor) {
	return &AnalyticV2Processor{
		repository: repository,
		analyticHandler: func(s *filter.StreamProcessorService, m *sarama.ConsumerMessage) {
			timberWolf, err := filter.ConvertKafkaMessageToTimberWolf(m)
			if err != nil {
				log.Warn(err.Error())
				return
			}

			if idx := getTimberWolfV2Index(timberWolf); idx != -1 {
				timberWolvesV2[idx].Counter++
			} else {
				timberWolvesV2 = append(timberWolvesV2, timberWolf)
			}
		},
	}
}

func (p *AnalyticV2Processor) Start() {
	log.Infof("Starting Analytic Processor...")
	interval := viper.GetString("WINDOWING_INTERVAL_IN_SECOND_V2") + "s"

	c := cron.New()
	c.AddFunc("@every "+interval, func() {
		p.saveToDatabase()
		log.Info(timberWolvesV2)
		timberWolvesV2 = []models.TimberWolf{}
	})
	go c.Start()
}

func (p *AnalyticV2Processor) Stop() {
	c.Stop()
}

func (p *AnalyticV2Processor) saveToDatabase() {
	for _, timberWolf := range timberWolvesV2 {
		p.repository.SaveV2(timberWolf)
	}

	log.Info("Logs aggregation stored!")
}

func (p *AnalyticV2Processor) GetHandler() processor.StreamHandler {
	return p.analyticHandler
}

func getTimberWolfV2Index(t models.TimberWolf) int {
	for idx, timberWolf := range timberWolvesV2 {
		if timberWolf.ApplicationName == t.ApplicationName && timberWolf.NodeId == t.NodeId && timberWolf.LogLevel == t.LogLevel {
			return idx
		}
	}

	return -1
}
