package processor

import (
	"github.com/BaritoLog/unclog-worker/filter"
	"github.com/Shopify/sarama"
)

type (
	StreamHandler = func(s *filter.StreamProcessorService, m *sarama.ConsumerMessage)

	Processor interface {
		GetHandler() StreamHandler
		Start()
		Stop()
	}

// 	Daemon interface {
// 		start()
// 		stop()
// 	}

// 	CompositeDaemon struct {
// 		daemon StreamHandler
// 		daemon Daemon
// 	}
// )

// func NewCompositeDaemon() {

// }

// func (p *CompositeDaemon) start() {
// 	handler.start
// 	daemon.start
// }

// func (p *CompositeDaemon) stop() {

// }
)
