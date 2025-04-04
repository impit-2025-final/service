package repository

import (
	"context"
	domain "service/internal/domain"

	"gorm.io/gorm"
)

type GormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

func (r *GormRepository) CreateNetworkTraffic(ctx context.Context, traffic *domain.NetworkTraffic) error {
	return r.db.Create(traffic).Error
}

func (r *GormRepository) CreateNetworkTrafficBatch(ctx context.Context, traffic []domain.NetworkTraffic) error {
	return r.db.Create(traffic).Error
}

func (r *GormRepository) CreateDockerInfo(ctx context.Context, dockerInfo domain.DockerInfo) error {
	return r.db.Create(dockerInfo).Error
}
