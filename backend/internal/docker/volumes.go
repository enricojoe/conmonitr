package docker

import (
	"context"

	"github.com/docker/docker/api/types/volume"
)

// VolumeSummary is the list-view DTO for Docker volumes sent to the frontend.
type VolumeSummary struct {
	Name       string `json:"name"`
	Driver     string `json:"driver"`
	Mountpoint string `json:"mountpoint"`
	Scope      string `json:"scope"`
	CreatedAt  string `json:"createdAt"`
}

// ListVolumes returns all volumes known to the Docker engine.
func (s *Service) ListVolumes(ctx context.Context) ([]VolumeSummary, error) {
	resp, err := s.cli.VolumeList(ctx, volume.ListOptions{})
	if err != nil {
		return nil, err
	}

	out := make([]VolumeSummary, 0, len(resp.Volumes))
	for _, v := range resp.Volumes {
		out = append(out, VolumeSummary{
			Name:       v.Name,
			Driver:     v.Driver,
			Mountpoint: v.Mountpoint,
			Scope:      v.Scope,
			CreatedAt:  v.CreatedAt,
		})
	}
	return out, nil
}
