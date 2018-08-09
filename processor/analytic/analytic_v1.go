package analytic

import (
	"fmt"

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
		},
	}
}

func (p *AnalyticV1Processor) Start() {
	log.Infof("Starting Analytic Processor...")
	interval := viper.GetString("WINDOWING_INTERVAL_IN_SECOND_V1") + "s"

	c := cron.New()
	c.AddFunc("@every "+interval, func() {
		p.saveToDatabase()
		resetAll(timberWolvesV1)
		fmt.Println("Every " + interval)
	})
	go c.Start()
}

func (p *AnalyticV1Processor) Stop() {
	c.Stop()
}

func (p *AnalyticV1Processor) saveToDatabase() {
	fmt.Println("Saving to database...")
	p.saveAllMetric()
	fmt.Println("Logs saving...")
	fmt.Println("Logs aggregation stored!")
}

func (p *AnalyticV1Processor) GetHandler() processor.StreamHandler {
	return p.analyticHandler
}

func (p *AnalyticV1Processor) saveAllMetric() {
	for _, timberWolf := range timberWolvesV1 {
		p.repository.SaveV1(timberWolf)
	}
}
func getTimberWolfV1Index(t models.TimberWolf) int {
	for idx, timberWolf := range timberWolvesV1 {
		if timberWolf.LogLevel == t.LogLevel {
			return idx
		}
	}

	return -1
}
