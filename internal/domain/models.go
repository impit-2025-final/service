package domain

import (
	"context"
	"encoding/json"

	"gorm.io/gorm"
)

type NodeInfo struct {
	gorm.Model
	Token    string   `json:"token" gorm:"uniqueIndex"`
	Hostname *string  `json:"hostname"`
	NodeName string   `json:"node_name"`
	IpsJSON  string   `json:"ips_json" gorm:"column:ips_json"`
	Ips      []string `json:"ips" gorm:"-"`
}

type DockerInfo struct {
	gorm.Model
	Containers []ContainerInfo `json:"containers" gorm:"foreignKey:DockerInfoID;references:ID"`
	Networks   []NetworkInfo   `json:"networks" gorm:"foreignKey:DockerInfoID;references:ID"`
	NodeInfoID uint            `json:"node_info_id" gorm:"index"`
}

type ContainerInfo struct {
	gorm.Model
	ContainerID       string            `json:"container_id"`
	ContainerName     string            `json:"container_name"`
	IP                string            `json:"ip"`
	Status            string            `json:"status"`
	LabelsJSON        string            `json:"labels_json" gorm:"column:labels_json"`
	Labels            map[string]string `json:"labels" gorm:"-"`
	AdditionalIPsJSON string            `json:"additional_ips_json" gorm:"column:additional_ips_json"`
	AdditionalIPs     []string          `json:"additional_ips" gorm:"-"`
	DockerInfoID      uint              `json:"docker_info_id" gorm:"index"`
}

type NetworkInfo struct {
	gorm.Model
	Name           string   `json:"name"`
	Subnet         string   `json:"subnet"`
	Gateway        string   `json:"gateway"`
	ContainersJSON string   `json:"containers_json" gorm:"column:containers_json"`
	Containers     []string `json:"containers" gorm:"-"`
	DockerInfoID   uint     `json:"docker_info_id" gorm:"index"`
}

func (n *NodeInfo) BeforeSave(tx *gorm.DB) (err error) {
	if n.Ips != nil {
		ipsJSON, err := json.Marshal(n.Ips)
		if err != nil {
			return err
		}
		n.IpsJSON = string(ipsJSON)
	}

	return nil
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
	CreateDockerInfo(ctx context.Context, dockerInfo DockerInfo, nodeId uint) error
	CreateNodeInfo(ctx context.Context, token string, nodeName string) (NodeInfo, error)
	UpdateNodeInfo(ctx context.Context, nodeInfo NodeInfo) error
	GetNodeInfo(ctx context.Context, token string) (NodeInfo, error)
}

type UseCase interface {
	UpdateNetworkTraffic(ctx context.Context, networkTraffic []NetworkTraffic) error
	UpdateDockerInfo(ctx context.Context, dockerInfo DockerInfo) error
	UpdateNodeInfo(ctx context.Context, nodeInfo NodeInfo) error
	GetNodeInfo(ctx context.Context, token string) (NodeInfo, error)
}
