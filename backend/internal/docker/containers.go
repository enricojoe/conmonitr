package docker

import (
	"context"
	"strings"

	"github.com/docker/docker/api/types/container"
)

// Port describes a published/exposed container port.
type Port struct {
	IP          string `json:"ip,omitempty"`
	PrivatePort uint16 `json:"privatePort"`
	PublicPort  uint16 `json:"publicPort,omitempty"`
	Type        string `json:"type"`
}

// ContainerSummary is the list-view DTO sent to the frontend.
type ContainerSummary struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Image   string `json:"image"`
	State   string `json:"state"`
	Status  string `json:"status"`
	Created int64  `json:"created"`
	Ports   []Port `json:"ports"`
}

// NetworkInfo summarises one attached network.
type NetworkInfo struct {
	Name       string `json:"name"`
	IPAddress  string `json:"ipAddress"`
	Gateway    string `json:"gateway"`
	MacAddress string `json:"macAddress"`
}

// MountInfo summarises one mount/volume.
type MountInfo struct {
	Type        string `json:"type"`
	Source      string `json:"source"`
	Destination string `json:"destination"`
	Name        string `json:"name,omitempty"`
	RW          bool   `json:"rw"`
}

// ContainerDetail is the inspect-view DTO sent to the frontend.
type ContainerDetail struct {
	ID            string        `json:"id"`
	Name          string        `json:"name"`
	Image         string        `json:"image"`
	State         string        `json:"state"`
	Status        string        `json:"status"`
	Created       string        `json:"created"`
	Command       []string      `json:"command"`
	Tty           bool          `json:"tty"`
	RestartPolicy string        `json:"restartPolicy"`
	Env           []string      `json:"env"`
	Ports         []Port        `json:"ports"`
	Networks      []NetworkInfo `json:"networks"`
	Mounts        []MountInfo   `json:"mounts"`
}

func cleanName(names []string) string {
	if len(names) == 0 {
		return ""
	}
	return strings.TrimPrefix(names[0], "/")
}

// List returns all containers (running and stopped).
func (s *Service) List(ctx context.Context) ([]ContainerSummary, error) {
	items, err := s.cli.ContainerList(ctx, container.ListOptions{All: true})
	if err != nil {
		return nil, err
	}

	out := make([]ContainerSummary, 0, len(items))
	for _, c := range items {
		ports := make([]Port, 0, len(c.Ports))
		for _, p := range c.Ports {
			ports = append(ports, Port{
				IP:          p.IP,
				PrivatePort: p.PrivatePort,
				PublicPort:  p.PublicPort,
				Type:        p.Type,
			})
		}
		out = append(out, ContainerSummary{
			ID:      c.ID,
			Name:    cleanName(c.Names),
			Image:   c.Image,
			State:   c.State,
			Status:  c.Status,
			Created: c.Created,
			Ports:   ports,
		})
	}
	return out, nil
}

// Inspect returns detailed information about a single container.
func (s *Service) Inspect(ctx context.Context, id string) (*ContainerDetail, error) {
	info, err := s.cli.ContainerInspect(ctx, id)
	if err != nil {
		return nil, err
	}

	d := &ContainerDetail{
		ID:      info.ID,
		Name:    strings.TrimPrefix(info.Name, "/"),
		Created: info.Created,
	}
	if info.Config != nil {
		d.Image = info.Config.Image
		d.Command = append(info.Config.Entrypoint, info.Config.Cmd...)
		d.Tty = info.Config.Tty
		d.Env = info.Config.Env
		for p := range info.Config.ExposedPorts {
			d.Ports = append(d.Ports, Port{
				PrivatePort: uint16(p.Int()),
				Type:        p.Proto(),
			})
		}
	}
	if info.State != nil {
		d.State = info.State.Status
		d.Status = info.State.Status
	}
	if info.HostConfig != nil {
		d.RestartPolicy = string(info.HostConfig.RestartPolicy.Name)
	}
	if info.NetworkSettings != nil {
		for name, n := range info.NetworkSettings.Networks {
			d.Networks = append(d.Networks, NetworkInfo{
				Name:       name,
				IPAddress:  n.IPAddress,
				Gateway:    n.Gateway,
				MacAddress: n.MacAddress,
			})
		}
	}
	for _, m := range info.Mounts {
		d.Mounts = append(d.Mounts, MountInfo{
			Type:        string(m.Type),
			Source:      m.Source,
			Destination: m.Destination,
			Name:        m.Name,
			RW:          m.RW,
		})
	}
	return d, nil
}

// Start starts a stopped container.
func (s *Service) Start(ctx context.Context, id string) error {
	return s.cli.ContainerStart(ctx, id, container.StartOptions{})
}

// Stop gracefully stops a running container.
func (s *Service) Stop(ctx context.Context, id string) error {
	return s.cli.ContainerStop(ctx, id, container.StopOptions{})
}

// Restart restarts a container.
func (s *Service) Restart(ctx context.Context, id string) error {
	return s.cli.ContainerRestart(ctx, id, container.StopOptions{})
}

// Remove deletes a container, optionally forcing removal of a running one.
func (s *Service) Remove(ctx context.Context, id string, force bool) error {
	return s.cli.ContainerRemove(ctx, id, container.RemoveOptions{Force: force})
}
