package domain

import (
	"context"
	"encoding/json"

	"gorm.io/gorm"
)

type DockerInfo struct {
	gorm.Model
	ID         string          `json:"id" gorm:"uniqueIndex"`
	Containers []ContainerInfo `json:"containers" gorm:"foreignKey:DockerInfoID;references:ID"`
	Networks   []NetworkInfo   `json:"networks" gorm:"foreignKey:DockerInfoID;references:ID"`
}

type ContainerInfo struct {
	gorm.Model
	ID                string            `json:"id" gorm:"uniqueIndex"`
	Name              string            `json:"name"`
	IP                string            `json:"ip"`
	Status            string            `json:"status"`
	LabelsJSON        string            `json:"labels_json" gorm:"column:labels_json"`
	Labels            map[string]string `json:"labels" gorm:"-"`
	AdditionalIPsJSON string            `json:"additional_ips_json" gorm:"column:additional_ips_json"`
	AdditionalIPs     []string          `json:"additional_ips" gorm:"-"`
	DockerInfoID      string            `json:"docker_info_id" gorm:"index"`
}

type NetworkInfo struct {
	gorm.Model
	Name           string   `json:"name"`
	Subnet         string   `json:"subnet"`
	Gateway        string   `json:"gateway"`
	ContainersJSON string   `json:"containers_json" gorm:"column:containers_json"`
	Containers     []string `json:"containers" gorm:"-"`
	DockerInfoID   string   `json:"docker_info_id" gorm:"index"`
}

func (c *ContainerInfo) BeforeSave(tx *gorm.DB) (err error) {
	if c.Labels != nil {
		labelsJSON, err := json.Marshal(c.Labels)
		if err != nil {
			return err
		}
		c.LabelsJSON = string(labelsJSON)
	}

	if c.AdditionalIPs != nil {
		ipsJSON, err := json.Marshal(c.AdditionalIPs)
		if err != nil {
			return err
		}
		c.AdditionalIPsJSON = string(ipsJSON)
	}

	return nil
}

func (c *ContainerInfo) AfterFind(tx *gorm.DB) (err error) {
	if c.LabelsJSON != "" {
		err = json.Unmarshal([]byte(c.LabelsJSON), &c.Labels)
		if err != nil {
			return err
		}
	}

	if c.AdditionalIPsJSON != "" {
		err = json.Unmarshal([]byte(c.AdditionalIPsJSON), &c.AdditionalIPs)
		if err != nil {
			return err
		}
	}

	return nil
}

func (n *NetworkInfo) BeforeSave(tx *gorm.DB) (err error) {
	if n.Containers != nil {
		containersJSON, err := json.Marshal(n.Containers)
		if err != nil {
			return err
		}
		n.ContainersJSON = string(containersJSON)
	}

	return nil
}

func (n *NetworkInfo) AfterFind(tx *gorm.DB) (err error) {
	if n.ContainersJSON != "" {
		err = json.Unmarshal([]byte(n.ContainersJSON), &n.Containers)
		if err != nil {
			return err
		}
	}

	return nil
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
