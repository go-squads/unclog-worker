package processor

import (
	"github.com/Shopify/sarama"
	"github.com/go-squads/unclog-worker/filter"
)

type (
	StreamHandler = func(s *filter.StreamProcessorService, m *sarama.ConsumerMessage)

	Processor interface {
		GetHandler() StreamHandler
		Start()
		Stop()
	}
)
