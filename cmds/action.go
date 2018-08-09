package cmds

import (
	"github.com/BaritoLog/go-boilerplate/srvkit"
	"github.com/Shopify/sarama"
	"github.com/go-squads/unclog-worker/filter"
	"github.com/go-squads/unclog-worker/processor/analytic"
	"github.com/go-squads/unclog-worker/processor/transformer"
	"github.com/urfave/cli"
)

//ActionStreamProcessorService .
func ActionStreamProcessorService(c *cli.Context) (err error) {
	brokers := configKafkaBrokers()
	groupID := configKafkaGroupId()
	topicSuffix := configKafkaTopicSuffix()
	newTopicEventName := configNewTopicEvent()
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	config.Version = sarama.V0_10_2_1
	factory := filter.NewKafkaFactory(brokers, config)

	transformer := transformer.NewTransformer()
	transformHandler := transformer.GetHandler()
	transformService := filter.NewStreamProcessorService(factory, groupID, topicSuffix, newTopicEventName, transformHandler)

	if err = transformService.Start(); err != nil {
		return
	}

	srvkit.GracefullShutdown(transformService.Close)
	return
}

//ActionStreamProcessorLogCountService .
func ActionStreamProcessorAllLogCounterService(c *cli.Context) (err error) {
	brokers := configKafkaBrokers()
	groupID := configKafkaGroupId()
	topicSuffix := "_processed"
	newTopicEventName := configNewTopicEvent()
	config := sarama.NewConfig()
	config.Version = sarama.V0_10_2_1
	factory := filter.NewKafkaFactory(brokers, config)

	analyticProcessor := analytic.NewAnalyticV1Processor(analytic.NewLogLevelRepositoryImpl())
	logCountHandler := analyticProcessor.GetHandler()
	logCountService := filter.NewStreamProcessorService(factory, groupID, topicSuffix, newTopicEventName, logCountHandler)

	analyticProcessor.Start()
	if err = logCountService.Start(); err != nil {
		return
	}

	srvkit.GracefullShutdown(logCountService.Close)
	srvkit.GracefullShutdown(analyticProcessor.Stop)
	return
}

func ActionStreamProcessorLogCounterByAppNameAndNodeIdServices(c *cli.Context) (err error) {
	brokers := configKafkaBrokers()
	groupID := configKafkaGroupId()
	topicSuffix := "_processed"
	newTopicEventName := configNewTopicEvent()
	config := sarama.NewConfig()
	config.Version = sarama.V0_10_2_1
	factory := filter.NewKafkaFactory(brokers, config)

	analyticProcessor := analytic.NewAnalyticV2Processor(analytic.NewLogLevelRepositoryImpl())
	logCountHandler := analyticProcessor.GetHandler()
	logCountService := filter.NewStreamProcessorService(factory, groupID, topicSuffix, newTopicEventName, logCountHandler)

	analyticProcessor.Start()
	if err = logCountService.Start(); err != nil {
		return
	}

	srvkit.GracefullShutdown(logCountService.Close)
	srvkit.GracefullShutdown(analyticProcessor.Stop)
	return
}
