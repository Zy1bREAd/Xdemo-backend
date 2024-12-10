package tmp

import "github.com/docker/go-connections/nat"

type ContainerBody struct {
	ID    string `json:"container_id"`
	Image string `json:"image"`
}

type ImageBody struct {
	Respository string `json:"respository"`
	Project     string `json:"project"`
	Name        string `json:"image" binding:"required"`
	Tag         string `json:"tag"`
	IsCover     bool   `json:"iscover"`
}

type ContainerCreateConfig struct {
	Image              string              `json:"image" binding:"required"`
	Name               string              `json:"container_name"`
	HostName           string              `json:"container_hostname"`
	ExposedPorts       map[nat.Port]string `json:"expose_ports"`
	PersistenceVolumes map[string]string   `json:"persistence_volumes"`
	TTY                bool                `json:"tty"`
	Interactive        bool                `json:"interactive"`
	Detach             bool                `json:"detach_run"`
	Env                []string            `json:"env_list"`
}
