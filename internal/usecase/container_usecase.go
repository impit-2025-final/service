package usecase

import (
	"context"
	domain "service/internal/domain"
)

type ContainerUseCase struct {
	repo domain.Repository
}

func NewContainerUseCase(repo domain.Repository) ContainerUseCase {
	return ContainerUseCase{repo: repo}
}

func (u *ContainerUseCase) UpdateNetworkTraffic(ctx context.Context, networkTraffic []domain.NetworkTraffic) error {
	return u.repo.CreateNetworkTrafficBatch(ctx, networkTraffic)
}

func (u *ContainerUseCase) UpdateDockerInfo(ctx context.Context, dockerInfo domain.DockerInfo) error {
	return u.repo.CreateDockerInfo(ctx, dockerInfo)
}
