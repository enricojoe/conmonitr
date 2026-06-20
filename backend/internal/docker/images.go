package docker

import (
	"context"

	"github.com/docker/docker/api/types/image"
)

// ImageSummary is the list-view DTO for Docker images sent to the frontend.
type ImageSummary struct {
	ID         string   `json:"id"`
	RepoTags   []string `json:"repoTags"`
	Size       int64    `json:"size"`
	Created    int64    `json:"created"`
	Containers int64    `json:"containers"`
}

// ListImages returns all non-intermediate images on the host.
func (s *Service) ListImages(ctx context.Context) ([]ImageSummary, error) {
	items, err := s.cli.ImageList(ctx, image.ListOptions{All: false})
	if err != nil {
		return nil, err
	}

	out := make([]ImageSummary, 0, len(items))
	for _, img := range items {
		containers := img.Containers
		if containers < 0 {
			containers = 0
		}
		out = append(out, ImageSummary{
			ID:         img.ID,
			RepoTags:   img.RepoTags,
			Size:       img.Size,
			Created:    img.Created,
			Containers: containers,
		})
	}
	return out, nil
}
