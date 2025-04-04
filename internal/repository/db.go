package repository

import (
	domain "service/internal/domain"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitDB(dbPath string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
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
