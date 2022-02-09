package mtjobrunner

import (
	types "github.com/Carlos-Descalzi/mtjobrunner/pkg/apis/messagelistener/v1alpha1"
	"go.uber.org/zap"
)

type RabbitMQSubscriber struct {
	name   string
	topic  string
	config types.RabbitMQListener
	logger *zap.SugaredLogger
}

func RabbitMQSubscriberNew(name string, topic string, config types.RabbitMQListener, logger *zap.SugaredLogger) *RabbitMQSubscriber {
	subscriber := RabbitMQSubscriber{name: name, topic: topic, logger: logger}
	subscriber.init()
	return &subscriber
}

func (s *RabbitMQSubscriber) init() {
	// TBC
}

func (h *RabbitMQSubscriber) Next(timeout uint32) (*Message, error) {
	return nil, nil
}

func (h *RabbitMQSubscriber) Stop() {
}
