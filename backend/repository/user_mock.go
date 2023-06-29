package repository

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"some-application/backend/model"
)

type MockUserRepository struct {
	Mock mock.Mock
}

func NewUserRepositoryMock() *MockUserRepository {
	return &MockUserRepository{}
}

func (m *MockUserRepository) InsertEmail(user *model.User) error {
	fmt.Printf("Value passed in: %v\n", user)
	args := m.Mock.Called(user)
	var r0 error
	if args.Get(0) != nil {
		r0 = args.Get(0).(error)
	}
	return r0
}

func (m *MockUserRepository) FindUserByEmail(email string) (*model.User, error) {
	fmt.Printf("Value passed in: %s\n", email)
	args := m.Mock.Called(email)
	var r0 *model.User
	if args.Get(0) != nil {
		r0 = args.Get(0).(*model.User)
	}
	var r1 error
	if args.Get(1) != nil {
		r1 = args.Get(1).(error)
	}
	return r0, r1
}

func (m *MockUserRepository) FindUserById(id uuid.UUID) (*model.User, error) {
	fmt.Printf("value passed in: %s\n", id.String())
	args := m.Mock.Called(id)
	var r0 *model.User
	if args.Get(0) != nil {
		r0 = args.Get(0).(*model.User)
	}
	var r1 error
	if args.Get(1) != nil {
		r1 = args.Get(1).(error)
	}
	return r0, r1
}
