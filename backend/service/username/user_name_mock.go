package username

import (
	"fmt"
	"github.com/stretchr/testify/mock"
	"some-application/backend/model"
	"some-application/backend/utils"
)

type MockService struct {
	Mock mock.Mock
}

func NewUsernameServiceMock() *MockService {
	return &MockService{}
}

func (m *MockService) AddUserName(username *model.UserName) *utils.Error {
	fmt.Printf("Value passed in: %v\n", username)
	args := m.Mock.Called(username)
	var r0 *utils.Error
	if args.Get(0) != nil {
		r0 = args.Get(0).(*utils.Error)
	}
	return r0
}
