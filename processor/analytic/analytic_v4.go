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
	AnalyticV4Processor struct {
		repository      LogLevelMetricRepository
		analyticHandler processor.StreamHandler
	}
)

var (
	timberWolvesV4 []models.TimberWolf
	cronJobV4      *gocron.Scheduler
)

func NewAnalyticV4Processor(repository LogLevelMetricRepository) (p *AnalyticV4Processor) {
	return &AnalyticV4Processor{
		repository: repository,
		analyticHandler: func(s *filter.StreamProcessorService, m *sarama.ConsumerMessage) {
			timberWolf, err := filter.ConvertKafkaMessageToTimberWolf(m)
			if err != nil {
				log.Warn(err.Error())
				return
			}

			if idx := getTimberWolfV4Index(timberWolf); idx != -1 {
				timberWolvesV4[idx].Counter++
			} else {
				timberWolvesV4 = append(timberWolvesV4, timberWolf)
			}
		},
	}
}

func (p *AnalyticV4Processor) Start() {
	log.Infof("Starting Analytic Processor...")
	interval := viper.GetInt("WINDOWING_INTERVAL_IN_SECOND_V2")

	cronJobV4 = gocron.NewScheduler()

	cronJobV4.Every(uint64(interval)).Seconds().Do(func() {
		p.saveToDatabase()
		log.Info(timberWolvesV4)
		timberWolvesV4 = []models.TimberWolf{}
	})
	cronJobV4.Start()
}

func (p *AnalyticV4Processor) Stop() {
	cronJobV4.Clear()
}

func (p *AnalyticV4Processor) saveToDatabase() {
	for _, timberWolf := range timberWolvesV4 {
		p.repository.SaveV4(timberWolf)
	}

	log.Info("Logs aggregation stored!")
}

func (p *AnalyticV4Processor) GetHandler() processor.StreamHandler {
	return p.analyticHandler
}

func getTimberWolfV4Index(t models.TimberWolf) int {
	for idx, timberWolf := range timberWolvesV4 {
		if timberWolf.ApplicationName == t.ApplicationName && timberWolf.NodeId == t.NodeId && timberWolf.LogLevel == t.LogLevel {
			return idx
		}
	}

	return -1
}

