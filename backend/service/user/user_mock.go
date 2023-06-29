package user

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"some-application/backend/model"
	"some-application/backend/utils"
)

type MockService struct {
	Mock mock.Mock
}

func NewUserServiceMock() *MockService {
	return &MockService{}
}

func (m *MockService) AddUser(user *model.User) *utils.Error {
	fmt.Printf("Value passed in: %v\n", user)
	args := m.Mock.Called(user)
	var r0 *utils.Error
	if args.Get(0) != nil {
		r0 = args.Get(0).(*utils.Error)
	}
	return r0
}

func (m *MockService) HashPassword(user *model.User) *utils.Error {
	fmt.Printf("Value passed in: %v\n", user)
	args := m.Mock.Called(user)
	var r0 *utils.Error
	if args.Get(0) != nil {
		r0 = args.Get(0).(*utils.Error)
	}
	return r0
}

func (m *MockService) UserAlreadyExists(newEmail string) (*utils.Error, bool) {
	fmt.Printf("Value passed in: %v\n", newEmail)
	args := m.Mock.Called(newEmail)
	var r0 *utils.Error
	if args.Get(0) != nil {
		r0 = args.Get(0).(*utils.Error)
	}
	var r1 bool
	if args.Get(1) != nil {
		r1 = args.Get(1).(bool)
	}
	return r0, r1
}

func (m *MockService) UserIdExists(id uuid.UUID) (*utils.Error, bool) {
	fmt.Printf("Value passed in: %v\n", id.String())
	args := m.Mock.Called(id)
	var r0 *utils.Error
	if args.Get(0) != nil {
		r0 = args.Get(0).(*utils.Error)
	}
	var r1 bool
	if args.Get(1) != nil {
		r1 = args.Get(1).(bool)
	}
	return r0, r1
}
