package kafka

import (
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log"
	"some-application/backend/kafka/message"
)

// https://www.tutorialsbuddy.com/write-data-to-a-kafka-topic-in-go-example

type ClientProducer interface {
	EmitEvent(eventType, topic, partitionKey string, payload message.NewUser) error
}

type kafkaProducer struct {
	Producer *kafka.Producer
}

func NewProducer() ClientProducer {
	producer := initProducer()
	return &kafkaProducer{
		Producer: producer,
	}
}

func initProducer() *kafka.Producer {
	producer, err := kafka.NewProducer(getProducerConfig())
	if err != nil {
		log.Fatal("Failed to start producer")
	}

	// delivery report handler for produced messages
	go func() {
		for e := range producer.Events() {
			switch event := e.(type) {
			case *kafka.Message:
				if event.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed %v\n", event.TopicPartition)
				} else {
					fmt.Printf("Delivered message %v\n", event.TopicPartition)
				}
			}
		}
	}()

	return producer
}

func getProducerConfig() *kafka.ConfigMap {
	return &kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
	}
}

func (p *kafkaProducer) EmitEvent(eventType, topic, partitionKey string, payload message.NewUser) error {
	headers := []kafka.Header{
		{
			Key:   "event-type",
			Value: []byte(eventType),
		},
	}

	encodedPayload, err := json.Marshal(payload)

	if err != nil {
		fmt.Println("error while marshaling payload")
		return err
	}

	msg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          encodedPayload,
		Headers:        headers,
		Key:            []byte(partitionKey),
	}

	return p.Producer.Produce(msg, nil)
}
