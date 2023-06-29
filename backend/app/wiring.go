package app

import (
	"log"
	"some-application/backend/database/postgresql"
	"some-application/backend/handler"
	"some-application/backend/kafka"
	"some-application/backend/repository"
	"some-application/backend/service/user"
	"some-application/backend/service/username"
)

type Wiring struct {
	HealthHandler      handler.HealthHandler
	CreateUserHandler  handler.CreateUserHandler
	AddUserNameHandler handler.AddUserNameHandler
}

func newWiring() Wiring {
	err := postgresql.Run()
	if err != nil {
		log.Fatalln(err)
	}

	gormDB, err := postgresql.New()
	if err != nil {
		log.Fatalln(err)
	}

	// kafka
	producer := kafka.NewProducer()

	// Repo
	userRepository := repository.NewUserRepository(gormDB)
	userNameRepository := repository.NewUserNameRepository(gormDB)

	// Service
	userService := user.NewUserService(userRepository)
	usernameService := username.NewUsernameService(userNameRepository)

	// Handler
	healthHandler := handler.NewHealthHandler()
	createUserHandler := handler.NewUserHandler(userService, producer)
	addUserNameHandler := handler.NewAddUserNameHandler(userService, usernameService, producer)

	return Wiring{
		HealthHandler:      healthHandler,
		CreateUserHandler:  createUserHandler,
		AddUserNameHandler: addUserNameHandler,
	}
}
