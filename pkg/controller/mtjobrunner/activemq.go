package mtjobrunner

import (
	"time"

	types "github.com/Carlos-Descalzi/mtjobrunner/pkg/apis/messagelistener/v1alpha1"
	"github.com/jjeffery/stomp"
	"go.uber.org/zap"
)

type ActiveMQSubscriber struct {
	active       bool
	destinations []string
	connection   *stomp.Conn
	config       types.ActiveMQListener
	channel      chan *Message
	logger       *zap.SugaredLogger
}

func ActiveMQSubscriberNew(name string, topics []string, config types.ActiveMQListener, logger *zap.SugaredLogger) *ActiveMQSubscriber {
	subscriber := ActiveMQSubscriber{
		destinations: topics,
		config:       config,
		logger:       logger,
		channel:      make(chan *Message),
	}
	subscriber.init()
	return &subscriber
}

func (s *ActiveMQSubscriber) Next(timeout time.Duration) (*Message, error) {
	select {
	case r := <-s.channel:
		return r, nil
	case <-time.After(timeout):
		return nil, nil
	}
}

func (s *ActiveMQSubscriber) Stop() {
	s.active = false
	if s.connection != nil {
		s.connection.Disconnect()
	}
	close(s.channel)
}

func (s *ActiveMQSubscriber) init() {

	conn, err := stomp.Dial(s.config.Network, s.config.Address)

	if err == nil {
		s.connection = conn
		s.active = true

		for i := 0; i < len(s.destinations); i++ {
			subscription, err := conn.Subscribe(s.destinations[i], stomp.AckAuto)

			if err == nil {
				go s.receive(s.destinations[i], subscription)
			} else {
				s.logger.Errorf("Unable to subscribe to %s", s.destinations[i])
			}
		}
	} else {
		s.logger.Errorf("Unable to connect to ActiveMQ %s %s", s.config.Network, s.config.Address)
	}
}

func (s *ActiveMQSubscriber) receive(destination string, subscription *stomp.Subscription) {
	for s.active {
		message, ok := <-subscription.C

		if ok {
			decoded := s.decodeMessage(message, destination)

			if decoded != nil {
				s.channel <- decoded
			}
		}
	}
}

func (s *ActiveMQSubscriber) decodeMessage(mqMessage *stomp.Message, destination string) *Message {

	message := Message{
		Properties: make(map[string]string),
		Payload:    mqMessage.Body,
		Topic:      destination,
	}

	message.Properties["ContentType"] = mqMessage.ContentType
	message.Properties["Destination"] = mqMessage.Destination

	if mqMessage.Header != nil {
		for i := 0; i < mqMessage.Header.Len(); i++ {
			key, value := mqMessage.Header.GetAt(i)
			message.Properties["Header_"+key] = value
		}
	}

	return &message
}
