package repository

import (
	"fmt"
	"service/internal/config"
	domain "service/internal/domain"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB(conf config.DBConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Europe/Moscow", conf.Host, conf.User, conf.Password, conf.DBName, conf.Port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(
		&domain.NetworkTraffic{},
		&domain.NetworkInfo{},
		&domain.ContainerInfo{},
		&domain.DockerInfo{},
	)
	if err != nil {
		return nil, err
	}

	return db, nil
}
