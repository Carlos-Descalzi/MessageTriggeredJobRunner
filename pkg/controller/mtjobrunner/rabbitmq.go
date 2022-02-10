package mtjobrunner

import (
	"strconv"
	"time"

	types "github.com/Carlos-Descalzi/mtjobrunner/pkg/apis/messagelistener/v1alpha1"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

type RabbitMQSubscriber struct {
	name       string
	active     bool
	topics     []string
	channel    chan *Message
	connection *amqp.Connection
	config     types.RabbitMQListener
	logger     *zap.SugaredLogger
}

func RabbitMQSubscriberNew(name string, topics []string, config types.RabbitMQListener, logger *zap.SugaredLogger) *RabbitMQSubscriber {
	subscriber := RabbitMQSubscriber{
		name:    name,
		topics:  topics,
		logger:  logger,
		channel: make(chan *Message),
	}
	subscriber.init()
	return &subscriber
}

func (s *RabbitMQSubscriber) init() {

	conn, err := amqp.Dial(s.config.Url)

	if err == nil {
		s.connection = conn
		channel, err := s.connection.Channel()

		if err == nil {
			s.active = true
			for i := 0; i < len(s.topics); i++ {
				delivery, err := channel.Consume(s.topics[i], "", true, false, false, false, amqp.Table{})

				if err == nil {
					go s.consume(s.topics[i], delivery)
				} else {
					s.logger.Errorf("Unable to consume %s: %s", s.topics[i], err)
				}
			}
		} else {
			s.logger.Error("Unable to open RabbitMQ channel ", err)
		}

	} else {
		s.logger.Errorf("Unable to connect to RabbitMQ %s: %s", s.config.Url, err)
	}
}

func (s *RabbitMQSubscriber) consume(topic string, delivery <-chan amqp.Delivery) {
	for s.active {
		inconmingMsg, ok := <-delivery

		if ok {
			message := s.decodeMessage(inconmingMsg, topic)
			if message != nil {
				s.channel <- message
			}
		}
	}
}

func (s *RabbitMQSubscriber) decodeMessage(inputMsg amqp.Delivery, topic string) *Message {
	message := Message{Payload: inputMsg.Body, Properties: make(map[string]string), Topic: topic}

	message.Properties["ContentType"] = inputMsg.ContentType
	message.Properties["ContentEncoding"] = inputMsg.ContentEncoding
	message.Properties["CorrelationId"] = inputMsg.CorrelationId
	message.Properties["Priority"] = strconv.Itoa(int(inputMsg.Priority))
	message.Properties["ReplyTo"] = inputMsg.ReplyTo
	message.Properties["MessageId"] = inputMsg.MessageId
	message.Properties["Type"] = inputMsg.Type
	message.Properties["UserId"] = inputMsg.UserId
	message.Properties["AppId"] = inputMsg.AppId
	message.Properties["Exchange"] = inputMsg.Exchange
	message.Properties["RoutingKey"] = inputMsg.RoutingKey

	return &message
}

func (s *RabbitMQSubscriber) Next(timeout time.Duration) (*Message, error) {
	select {
	case r := <-s.channel:
		return r, nil
	case <-time.After(timeout):
		return nil, nil
	}
}

func (s *RabbitMQSubscriber) Stop() {
	s.active = false
	s.connection.Close()
}
