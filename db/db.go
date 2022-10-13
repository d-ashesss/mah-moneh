package db

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect(host, user, password, dbname string) (err error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s database=%s", host, user, password, dbname)
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: logger.Discard})
	return err
}
