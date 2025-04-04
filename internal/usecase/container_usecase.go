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

func (u *ContainerUseCase) UpdateDockerInfo(ctx context.Context, dockerInfo domain.DockerInfo, nodeId uint) error {
	return u.repo.CreateDockerInfo(ctx, dockerInfo, nodeId)
}

func (u *ContainerUseCase) UpdateNodeInfo(ctx context.Context, nodeInfo domain.NodeInfo) error {
	return u.repo.UpdateNodeInfo(ctx, nodeInfo)
}

func (u *ContainerUseCase) CreateNodeInfo(ctx context.Context, token string, nodeName string) (domain.NodeInfo, error) {
	return u.repo.CreateNodeInfo(ctx, token, nodeName)
}

func (u *ContainerUseCase) GetNodeInfo(ctx context.Context, token string) (domain.NodeInfo, error) {
	return u.repo.GetNodeInfo(ctx, token)
}
