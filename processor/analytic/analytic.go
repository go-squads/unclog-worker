package analytic

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/go-squads/unclog-worker/filter"
	"github.com/go-squads/unclog-worker/models"
	"github.com/go-squads/unclog-worker/processor"
	"github.com/prometheus/common/log"
	"github.com/robfig/cron"
)

type (
	AnalyticProcessor struct {
		repository      LogLevelMetricRepository
		analyticHandler processor.StreamHandler
	}
)

var (
	c            cron.Cron
	timberWolves []models.TimberWolf
)

func NewAnalyticProcessor(repository LogLevelMetricRepository) (p *AnalyticProcessor) {
	return &AnalyticProcessor{
		repository: repository,
		analyticHandler: func(s *filter.StreamProcessorService, m *sarama.ConsumerMessage) {

			timberWolf, err := filter.ConvertKafkaMessageToTimberWolf(m)
			if err != nil {
				log.Warn(err.Error())
				return
			}

			if idx := getTimberWolfIndex(timberWolf); idx != -1 {
				timberWolves[idx].Counter++
			} else {
				timberWolves = append(timberWolves, timberWolf)
			}
		},
	}
}

func (p *AnalyticProcessor) Start() {
	log.Infof("Starting Analytic Processor...")

	c := cron.New()
	c.AddFunc("@every 1m", func() {
		p.saveToDatabase()
		fmt.Println("Every 1m")
	})
	go c.Start()
}

func (p *AnalyticProcessor) Stop() {
	c.Stop()
}

func (p *AnalyticProcessor) saveToDatabase() {
	fmt.Println("Saving to database...")
	p.saveAllMetric()
	fmt.Print("\nLogs saving...\n")
	p.resetAll()
	fmt.Println("Database work done, logs reset!")
}

func (p *AnalyticProcessor) GetHandler() processor.StreamHandler {
	return p.analyticHandler
}

func (p *AnalyticProcessor) saveAllMetric() {
	for _, timberWolf := range timberWolves {
		p.repository.Save(timberWolf)
	}
}

func getTimberWolfIndex(t models.TimberWolf) int {
	for idx, timberWolf := range timberWolves {
		if timberWolf.ApplicationName == t.ApplicationName && timberWolf.NodeId == t.NodeId && timberWolf.LogLevel == t.LogLevel {
			return idx
		}
	}

	return -1
}

func (p *AnalyticProcessor) resetAll() {
	timberWolves = []models.TimberWolf{}
}
