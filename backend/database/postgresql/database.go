package postgresql

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func New() (*gorm.DB, error) {
	url := fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=%v",
		DbHost, DbPort, DbUser, DbPassword, DbName, DbSslMode)

	db, err := gorm.Open(postgres.Open(url), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
