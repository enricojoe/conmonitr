package docker

import (
	"context"

	"github.com/docker/docker/api/types/network"
)

// NetworkSummary is the list-view DTO for Docker networks sent to the frontend.
type NetworkSummary struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Driver   string `json:"driver"`
	Scope    string `json:"scope"`
	Created  string `json:"created"`
	Internal bool   `json:"internal"`
}

// ListNetworks returns all networks known to the Docker engine.
func (s *Service) ListNetworks(ctx context.Context) ([]NetworkSummary, error) {
	items, err := s.cli.NetworkList(ctx, network.ListOptions{})
	if err != nil {
		return nil, err
	}

	out := make([]NetworkSummary, 0, len(items))
	for _, n := range items {
		out = append(out, NetworkSummary{
			ID:       n.ID,
			Name:     n.Name,
			Driver:   n.Driver,
			Scope:    n.Scope,
			Created:  n.Created.Format("2006-01-02T15:04:05Z07:00"),
			Internal: n.Internal,
		})
	}
	return out, nil
}
