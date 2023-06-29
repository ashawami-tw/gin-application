package kafka

import (
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log"
	"some-application/backend/kafka/message"
)

type ClientConsumer interface {
	Start() error
	Stop() error
}

type kafkaConsumer struct {
	Consumer *kafka.Consumer
}

func NewConsumer(topic []string) *kafkaConsumer {
	consumer := initConsumer(topic)
	return &kafkaConsumer{
		Consumer: consumer,
	}
}

func initConsumer(topic []string) *kafka.Consumer {
	consumer, err := kafka.NewConsumer(getConsumerConfig())
	if err != nil {
		log.Fatal("Failed to start consumer")
	}

	err = consumer.SubscribeTopics(topic, nil)
	if err != nil {
		log.Fatal("failed to subscribe topics")
	}
	return consumer
}

func getConsumerConfig() *kafka.ConfigMap {
	return &kafka.ConfigMap{
		"bootstrap.servers":  "localhost:9092",
		"session.timeout.ms": 6000,
		"auto.offset.reset":  "earliest",
		"group.id":           "new-user",
	}
}

func (c *kafkaConsumer) Start() error {
	run := true
	for run {
		// -1 - indefinite wait
		msg, err := c.Consumer.ReadMessage(-1)
		if err != nil {
			if err.(kafka.Error).Code() == kafka.ErrTransport {
				continue
			} else {
				return err
			}
		} else {
			var data message.NewUser
			err = json.Unmarshal(msg.Value, &data)
			if err != nil {
				fmt.Println("error while json unmarshal kafka consumed message")
			} else {
				fmt.Printf("Message on %s: %v\n", msg.TopicPartition, data)
				if msg.Headers != nil {
					fmt.Printf("%% Headers: %v\n", msg.Headers)
				}
			}
		}
	}
	return nil
}

func (c *kafkaConsumer) Stop() error {
	err := c.Consumer.Close()
	if err != nil {
		return err
	}
	return nil
}
