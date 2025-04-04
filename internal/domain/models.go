package domain

import (
	"context"

	"gorm.io/gorm"
)

type DockerInfo struct {
	gorm.Model
	ID         string          `json:"id" gorm:"uniqueIndex"`
	Containers []ContainerInfo `json:"containers"`
	Networks   []NetworkInfo   `json:"networks"`
}

type ContainerInfo struct {
	gorm.Model
	ID            string            `json:"id" gorm:"uniqueIndex"`
	Name          string            `json:"name"`
	IP            string            `json:"ip"`
	Status        string            `json:"status"`
	Labels        map[string]string `json:"labels" gorm:"-"`
	AdditionalIPs []string          `json:"additional_ips" gorm:"-"`
	DockerInfoID  string            `json:"docker_info_id" gorm:"uniqueIndex"`
}

type NetworkInfo struct {
	gorm.Model
	Name         string   `json:"name"`
	Subnet       string   `json:"subnet"`
	Gateway      string   `json:"gateway"`
	Containers   []string `json:"containers" gorm:"-"`
	DockerInfoID string   `json:"docker_info_id" gorm:"uniqueIndex"`
}

type NetworkTraffic struct {
	gorm.Model
	SourceIP      string `json:"source_ip"`
	DestinationIP string `json:"destination_ip"`
	Protocol      string `json:"protocol"`
	Interface     string `json:"interface"`
	Bytes         int64  `json:"bytes"`
	Packets       int64  `json:"packets"`
	ContainerID   string `json:"container_id"`
	ContainerName string `json:"container_name"`
	SrcPort       uint16 `json:"src_port"`
	DstPort       uint16 `json:"dst_port"`
}

type Repository interface {
	CreateNetworkTraffic(ctx context.Context, traffic *NetworkTraffic) error
	CreateNetworkTrafficBatch(ctx context.Context, traffic []NetworkTraffic) error

	CreateDockerInfo(ctx context.Context, dockerInfo DockerInfo) error
}

type UseCase interface {
	UpdateNetworkTraffic(ctx context.Context, networkTraffic []NetworkTraffic) error
	UpdateDockerInfo(ctx context.Context, dockerInfo DockerInfo) error
}
