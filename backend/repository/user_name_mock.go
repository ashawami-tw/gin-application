package repository

import (
	"fmt"
	"github.com/stretchr/testify/mock"
	"some-application/backend/model"
)

type MockUsernameRepository struct {
	Mock mock.Mock
}

func NewUsernameRepositoryMock() *MockUsernameRepository {
	return &MockUsernameRepository{}
}

func (m *MockUsernameRepository) InsertName(username *model.UserName) error {
	fmt.Printf("Value passed in: %v\n", username)
	args := m.Mock.Called(username)
	var r0 error
	if args.Get(0) != nil {
		r0 = args.Get(0).(error)
	}
	return r0
}
