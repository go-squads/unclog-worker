package analytic

import (
	"fmt"
	"sync/atomic"

	"github.com/BaritoLog/unclog-worker/filter"
	"github.com/BaritoLog/unclog-worker/processor"
	"github.com/Shopify/sarama"
	"github.com/prometheus/common/log"
	"github.com/robfig/cron"
)

type (
	AnalyticProcessor struct {
		repository      LogLevelMetricRepository
		unlistedMetric  uint32
		errorMetric     uint32
		infoMetric      uint32
		traceMetric     uint32
		debugMetric     uint32
		warnMetric      uint32
		analyticHandler processor.StreamHandler
	}
)

var (
	isStart bool = false
	c       cron.Cron
)

func NewAnalyticProcessor(repository LogLevelMetricRepository) (p *AnalyticProcessor) {
	return &AnalyticProcessor{
		// repository
		repository: repository,
		analyticHandler: func(s *filter.StreamProcessorService, m *sarama.ConsumerMessage) {

			timberWolf, err := filter.ConvertKafkaMessageToTimberWolf(m)
			if err != nil {
				log.Warn(err.Error())
				return
			}

			//unlisted = 0, error = 1, info = 2, trace = 3, debug = 4, warn = 5
			switch logLevelId := timberWolf.LogLevelId(); logLevelId {
			case 1:
				p.incrementError()
			case 2:
				p.incrementInfo()
			case 3:
				p.incrementTrace()
			case 4:
				p.incrementDebug()
			case 5:
				p.incrementWarn()
			default:
				p.incrementUnlisted()
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

	unlistedCount := atomic.LoadUint32(&p.unlistedMetric)
	p.saveMetric("UNLISTED", 0, unlistedCount)

	errorCount := atomic.LoadUint32(&p.errorMetric)
	p.saveMetric("ERROR", 1, errorCount)

	infoCount := atomic.LoadUint32(&p.infoMetric)
	p.saveMetric("INFO", 2, infoCount)

	traceCount := atomic.LoadUint32(&p.traceMetric)
	p.saveMetric("TRACE", 3, traceCount)

	debugCount := atomic.LoadUint32(&p.debugMetric)
	p.saveMetric("DEBUG", 4, debugCount)

	warnCount := atomic.LoadUint32(&p.warnMetric)
	p.saveMetric("WARN", 5, warnCount)
}

func (p *AnalyticProcessor) saveMetric(logLevel string, logLevelId int, count uint32) {
	metric := LogLevelMetric{
		logLevel:   logLevel,
		logLevelId: logLevelId,
		count:      count,
	}
	p.repository.Save(metric)
}

func (p *AnalyticProcessor) incrementUnlisted() {
	atomic.AddUint32(&p.unlistedMetric, 1)
}

func (p *AnalyticProcessor) incrementError() {
	atomic.AddUint32(&p.errorMetric, 1)
}

func (p *AnalyticProcessor) incrementInfo() {
	atomic.AddUint32(&p.infoMetric, 1)
}

func (p *AnalyticProcessor) incrementTrace() {
	atomic.AddUint32(&p.traceMetric, 1)
}

func (p *AnalyticProcessor) incrementDebug() {
	atomic.AddUint32(&p.debugMetric, 1)
}

func (p *AnalyticProcessor) incrementWarn() {
	atomic.AddUint32(&p.warnMetric, 1)
}

func (p *AnalyticProcessor) resetAll() {
	atomic.StoreUint32(&p.unlistedMetric, 0)
	atomic.StoreUint32(&p.errorMetric, 0)
	atomic.StoreUint32(&p.infoMetric, 0)
	atomic.StoreUint32(&p.traceMetric, 0)
	atomic.StoreUint32(&p.debugMetric, 0)
	atomic.StoreUint32(&p.warnMetric, 0)
}
