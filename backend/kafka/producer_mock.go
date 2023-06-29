package kafka

import (
	"fmt"
	"github.com/stretchr/testify/mock"
	"some-application/backend/kafka/message"
)

type MockProducer struct {
	Mock mock.Mock
}

func NewMockProducer() *MockProducer {
	return &MockProducer{}
}

func (m *MockProducer) EmitEvent(eventType, topic, partitionKey string, payload message.NewUser) error {
	fmt.Printf("Value passed in eventType: %v, topic: %v, partitionKey: %v, payload: %v\n", eventType, topic, partitionKey, payload)
	args := m.Mock.Called(eventType, topic, partitionKey, payload)

	var r0 error
	if args.Get(0) != nil {
		r0 = args.Get(0).(error)
	}
	return r0
}
