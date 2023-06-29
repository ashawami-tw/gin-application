package main

import (
	"fmt"
	"some-application/backend/kafka"
)

func main() {
	topic := []string{"user"}
	consumer := kafka.NewConsumer(topic)
	err := consumer.Start()
	if err != nil {
		err = consumer.Stop()
		if err != nil {
			fmt.Println(err)
		}
	}
}
