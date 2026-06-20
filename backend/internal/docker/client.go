// Package docker wraps the official Docker Engine SDK with the small surface
// ConMonitr needs: listing/inspecting containers, lifecycle actions, and
// streaming live stats and logs over the engine's Unix socket.
package docker

import (
	"context"

	"github.com/docker/docker/client"
)

// Service is a thin wrapper around the Docker SDK client.
type Service struct {
	cli *client.Client
}

// NewService creates a Docker client using the standard environment
// (DOCKER_HOST, etc.) with API-version negotiation so it stays compatible
// with whatever engine version is running locally.
func NewService() (*Service, error) {
	cli, err := client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)
	if err != nil {
		return nil, err
	}
	return &Service{cli: cli}, nil
}

// Ping verifies the engine is reachable and returns its API version string.
func (s *Service) Ping(ctx context.Context) (string, error) {
	p, err := s.cli.Ping(ctx)
	if err != nil {
		return "", err
	}
	return p.APIVersion, nil
}

// Close releases the underlying client.
func (s *Service) Close() error {
	return s.cli.Close()
}
