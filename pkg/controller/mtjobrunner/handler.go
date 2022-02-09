package mtjobrunner

import (
	"container/list"
	"fmt"

	ltypes "github.com/Carlos-Descalzi/mtjobrunner/pkg/apis/messagelistener/v1alpha1"
	jtypes "github.com/Carlos-Descalzi/mtjobrunner/pkg/apis/mtjob/v1alpha1"
	"go.uber.org/zap"
)

type subscriberError struct {
	message string
}

func (s subscriberError) Error() string {
	return s.message
}

type SubscriberHandler struct {
	namespace  string
	listeners  *list.List
	active     bool
	config     ltypes.MessageListener
	subscriber Subscriber
	logger     *zap.SugaredLogger
}

func SubscriberHandlerNew(namespace string, config ltypes.MessageListener, logger *zap.SugaredLogger) *SubscriberHandler {
	var listener = SubscriberHandler{
		namespace: namespace,
		listeners: list.New(),
		config:    config,
		logger:    logger,
	}
	return &listener
}

func (l *SubscriberHandler) AddListener(listener interface{}) {
	l.listeners.PushBack(listener)
}

func (l *SubscriberHandler) Start() error {

	if l.config.Spec.Kafka.IsSet() {
		l.subscriber = KafkaSubscriberNew(l.config.Name, l.config.Spec.Topic, l.config.Spec.Kafka, l.logger)
	} else if l.config.Spec.RabbitMQ.IsSet() {
		l.subscriber = RabbitMQSubscriberNew(l.config.Name, l.config.Spec.Topic, l.config.Spec.RabbitMQ, l.logger)
	} else {
		return subscriberError{
			fmt.Sprintf(
				"Unable to start subscriber handler %s/%s, listener configuration has no kafka or rabbitmq configuration",
				l.config.Namespace,
				l.config.Name,
			),
		}
	}

	l.active = true
	go l.loop()

	return nil
}

func (l *SubscriberHandler) loop() {
	for l.active {

		message, err := l.subscriber.Next(1000)

		if err == nil && message != nil {

			l.logger.Debugf(
				"Received message from topic %s, on listener %s",
				l.config.Spec.Topic,
				l.config.Name,
			)

			for listener := l.listeners.Front(); listener != nil; listener = listener.Next() {
				listener.Value.(Listener).MessageReceived(
					l.namespace,
					jtypes.Trigger{ListenerName: l.config.Name, Topic: l.config.Spec.Topic},
					message,
				)
			}
		}
	}
}

func (l *SubscriberHandler) Stop() {
	l.active = false
	if l.subscriber != nil {
		l.subscriber.Stop()
	}
}
