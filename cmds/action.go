package cmds

import (
	"github.com/BaritoLog/go-boilerplate/srvkit"
	"github.com/BaritoLog/unclog/filter"
	"github.com/Shopify/sarama"
	"github.com/prometheus/common/log"
	"github.com/urfave/cli"
)

//ActionStreamProcessorService .
func ActionStreamProcessorService(c *cli.Context) (err error) {
	brokers := configKafkaBrokers()
	groupID := configKafkaGroupId()
	//esUrl := configElasticsearchUrl()       //remove
	topicSuffix := configKafkaTopicSuffix() //_logs
	//topicSuffixProcessed := configKafkaTopicSuffix() //change this
	newTopicEventName := configNewTopicEvent()
	//newTopicEventNameProcessed := configNewTopicEvent() //changethis
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	// config.Producer.Retry.Max = maxRetry
	config.Producer.Return.Successes = true
	config.Version = sarama.V0_10_2_1

	factory := filter.NewKafkaFactory(brokers, config)

	transformHandler := func(s *filter.StreamProcessorService, m *sarama.ConsumerMessage) {

		timberWolf, err := filter.ConvertKafkaMessageToTimberWolf(m)
		if err != nil {
			log.Warn(err.Error())
			return
		}

		err = s.SendLogs("manan_logs_processed", timberWolf)
		if err != nil {
			log.Info(err)
			return
		}

		s.LogTimberWolf(timberWolf)
	}

	transformService := filter.NewStreamProcessorService(factory, groupID, topicSuffix, newTopicEventName, transformHandler)

	if err = transformService.Start(); err != nil {
		return
	}

	srvkit.GracefullShutdown(transformService.Close)

	return
}
