package filter

import (
	"strings"

	"github.com/BaritoLog/go-boilerplate/errkit"
	"github.com/Shopify/sarama"
	log "github.com/sirupsen/logrus"
)

const (
	ErrConvertKafkaMessage   = errkit.Error("Convert KafkaMessage Failed")
	ErrStore                 = errkit.Error("Store Failed")
	ErrConsumerWorker        = errkit.Error("Consumer Worker Failed")
	ErrMakeKafkaAdmin        = errkit.Error("Make kafka admin failed")
	ErrMakeNewTopicWorker    = errkit.Error("Make new topic worker failed")
	ErrSpawnWorkerOnNewTopic = errkit.Error("Spawn worker on new topic failed")
	ErrSpawnWorker           = errkit.Error("Span worker failed")
	ErrMakeSyncProducer      = errkit.Error("Make sync producer failed")
)

type StreamProcessorServiceInterface interface {
	Start() error
	Close()
	WorkerMap() map[string]ConsumerWorker
	NewTopicEventWorker() ConsumerWorker
	SendLogs(topic string, timberWolf TimberWolf) error
}

type SaramaMessageHandler = func(*StreamProcessorService, *sarama.ConsumerMessage)

type StreamProcessorService struct {
	factory           KafkaFactory
	groupID           string
	topicSuffix       string
	newTopicEventName string

	workerMap           map[string]ConsumerWorker
	admin               KafkaAdmin
	newTopicEventWorker ConsumerWorker

	lastError      error
	lastTimberWolf TimberWolf
	lastNewTopic   string
	handler        SaramaMessageHandler
}

func NewStreamProcessorService(factory KafkaFactory, groupID, topicSuffix, newTopicEventName string, handler SaramaMessageHandler) StreamProcessorServiceInterface {

	return &StreamProcessorService{
		factory:           factory,
		groupID:           groupID,
		topicSuffix:       topicSuffix,
		newTopicEventName: newTopicEventName,
		workerMap:         make(map[string]ConsumerWorker),
		handler:           handler,
	}
}

func (s *StreamProcessorService) SendLogs(topic string, timberWolf TimberWolf) (err error) {
	producer, err := s.factory.MakeSyncProducer()
	if err != nil {
		err = errkit.Concat(ErrMakeSyncProducer, err)
		return
	}

	message := ConvertTimberWolfToKafkaMessage(timberWolf, topic)
	_, _, err = producer.SendMessage(message)
	return
}

func (s *StreamProcessorService) Start() (err error) {

	admin, err := s.initAdmin()
	if err != nil {
		return errkit.Concat(ErrMakeKafkaAdmin, err)
	}

	worker, err := s.initNewTopicWorker()
	if err != nil {
		return errkit.Concat(ErrMakeNewTopicWorker, err)
	}

	worker.Start()

	for _, topic := range admin.Topics() {
		if strings.HasSuffix(topic, s.topicSuffix) {
			err := s.spawnLogsWorker(topic, sarama.OffsetNewest)
			if err != nil {
				s.logError(errkit.Concat(ErrSpawnWorker, err))
			}
		}
	}

	return
}

func (s *StreamProcessorService) initAdmin() (admin KafkaAdmin, err error) {
	admin, err = s.factory.MakeKafkaAdmin()
	s.admin = admin
	return
}

func (s *StreamProcessorService) initNewTopicWorker() (worker ConsumerWorker, err error) { // TODO: return worker
	topic := s.newTopicEventName
	consumer, err := s.factory.MakeClusterConsumer(s.groupID, topic, sarama.OffsetNewest)
	if err != nil {
		return
	}

	worker = NewConsumerWorker(topic, consumer)
	worker.OnSuccess(s.onNewTopicEvent)
	worker.OnError(s.logError)

	s.newTopicEventWorker = worker
	return
}

func (s StreamProcessorService) Close() {
	for _, worker := range s.workerMap {
		worker.Stop()
	}

	if s.admin != nil {
		s.admin.Close()
	}

	if s.newTopicEventWorker != nil {
		s.newTopicEventWorker.Stop()
	}
}

func (s *StreamProcessorService) getConsumerHandler() func(*sarama.ConsumerMessage) {
	return func(m *sarama.ConsumerMessage) {
		s.handler(s, m)
	}
}

func (s *StreamProcessorService) spawnLogsWorker(topic string, initialOffset int64) (err error) {

	consumer, err := s.factory.MakeClusterConsumer(s.groupID, topic, initialOffset)
	if err != nil {
		err = errkit.Concat(ErrConsumerWorker, err)
		return
	}

	worker := NewConsumerWorker(topic, consumer)
	worker.OnError(s.logError)
	worker.OnSuccess(s.getConsumerHandler())
	worker.Start()

	s.workerMap[topic] = worker

	return
}

func (s *StreamProcessorService) logError(err error) {
	s.lastError = err
	log.Warn(err.Error())
}

func (s *StreamProcessorService) LogTimberWolf(timberWolf TimberWolf) {
	s.lastTimberWolf = timberWolf
	log.Infof("TimberWolf: %v", timberWolf)
}

func (s *StreamProcessorService) logNewTopic(topic string) {
	s.lastNewTopic = topic
	log.Infof("New topic: %s", topic)
}

func (s *StreamProcessorService) onNewTopicEvent(message *sarama.ConsumerMessage) {
	topic := string(message.Value)

	err := s.spawnLogsWorker(topic, sarama.OffsetOldest)

	if err != nil {
		s.logError(errkit.Concat(ErrSpawnWorkerOnNewTopic, err))
		return
	}

	s.logNewTopic(topic)
}

func (s *StreamProcessorService) WorkerMap() map[string]ConsumerWorker {
	return s.workerMap
}

func (s *StreamProcessorService) NewTopicEventWorker() ConsumerWorker {
	return s.newTopicEventWorker
}
