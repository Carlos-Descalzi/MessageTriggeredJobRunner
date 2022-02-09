package mtjobrunner

import (
	"errors"
	"strings"
	"time"

	types "github.com/Carlos-Descalzi/mtjobrunner/pkg/apis/messagelistener/v1alpha1"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"go.uber.org/zap"
)

type KafkaSubscriber struct {
	name     string
	topic    string
	config   types.KafkaListener
	consumer *kafka.Consumer
	logger   *zap.SugaredLogger
}

func KafkaSubscriberNew(name string, topic string, config types.KafkaListener, logger *zap.SugaredLogger) *KafkaSubscriber {
	subscriber := KafkaSubscriber{topic: topic, config: config, logger: logger}
	subscriber.init()
	return &subscriber
}

func (h *KafkaSubscriber) init() {

	var config = kafka.ConfigMap{}

	for key, value := range h.config.Config {
		config.SetKey(key, kafka.ConfigValue(value))
	}

	consumer, err := kafka.NewConsumer(&config)

	if err != nil {
		panic(err)
	}
	h.consumer = consumer
	h.consumer.SubscribeTopics([]string{h.topic}, nil)
}

func (h *KafkaSubscriber) Next(timeout uint32) (*Message, error) {

	if h.consumer != nil {
		msg, err := h.consumer.ReadMessage(time.Duration(timeout * uint32(time.Millisecond)))

		if err != nil {
			return nil, err
		}
		return h.decodeMessage(msg), nil
	}

	return nil, errors.New("Kafka consumer not initialized or closed")
}

func (h KafkaSubscriber) decodeMessage(message *kafka.Message) *Message {
	var msg = Message{Properties: make(map[string]string)}
	msg.Properties["KEY"] = string(message.Key)
	for i := 0; i < len(message.Headers); i++ {
		msg.Properties[strings.ToUpper(message.Headers[i].Key)] = string(message.Headers[i].Value)
	}
	msg.Payload = message.Value

	return &msg
}

func (h *KafkaSubscriber) Stop() {
	if h.consumer != nil {
		h.consumer.Close()
		h.consumer = nil
	}
}
