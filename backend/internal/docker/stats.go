package docker

import (
	"context"
	"encoding/json"
	"io"
	"strings"
	"time"

	"github.com/docker/docker/api/types/container"
)

// Metric is a single computed performance snapshot for one container.
type Metric struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	Timestamp  int64   `json:"timestamp"` // unix milliseconds
	CPUPercent float64 `json:"cpuPercent"`
	MemUsage   uint64  `json:"memUsage"`
	MemLimit   uint64  `json:"memLimit"`
	MemPercent float64 `json:"memPercent"`
	NetRx      uint64  `json:"netRx"`
	NetTx      uint64  `json:"netTx"`
	BlkRead    uint64  `json:"blkRead"`
	BlkWrite   uint64  `json:"blkWrite"`
}

// StreamStats opens a streaming stats reader for one container and sends a
// computed Metric to out for every frame the engine emits (~1/sec). It blocks
// until ctx is cancelled, the stream ends, or a fatal error occurs.
func (s *Service) StreamStats(ctx context.Context, id, name string, out chan<- Metric) error {
	resp, err := s.cli.ContainerStats(ctx, id, true)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)
	for {
		var raw container.StatsResponse
		if err := dec.Decode(&raw); err != nil {
			if err == io.EOF || ctx.Err() != nil {
				return nil
			}
			return err
		}
		m := computeMetric(id, name, &raw)
		select {
		case out <- m:
		case <-ctx.Done():
			return nil
		}
	}
}

// computeMetric turns a raw engine stats frame into the values the dashboard
// displays, applying the standard `docker stats` CPU% formula.
func computeMetric(id, name string, raw *container.StatsResponse) Metric {
	m := Metric{
		ID:        id,
		Name:      name,
		Timestamp: time.Now().UnixMilli(),
	}

	// CPU percentage from usage deltas vs system-wide CPU time.
	cpuDelta := float64(raw.CPUStats.CPUUsage.TotalUsage) - float64(raw.PreCPUStats.CPUUsage.TotalUsage)
	sysDelta := float64(raw.CPUStats.SystemUsage) - float64(raw.PreCPUStats.SystemUsage)
	onlineCPUs := float64(raw.CPUStats.OnlineCPUs)
	if onlineCPUs == 0 {
		onlineCPUs = float64(len(raw.CPUStats.CPUUsage.PercpuUsage))
	}
	if sysDelta > 0 && cpuDelta > 0 {
		m.CPUPercent = (cpuDelta / sysDelta) * onlineCPUs * 100.0
	}

	// Memory: subtract page cache so the figure matches `docker stats`.
	usage := raw.MemoryStats.Usage
	if cache, ok := raw.MemoryStats.Stats["inactive_file"]; ok {
		if cache < usage {
			usage -= cache
		}
	} else if cache, ok := raw.MemoryStats.Stats["cache"]; ok {
		if cache < usage {
			usage -= cache
		}
	}
	m.MemUsage = usage
	m.MemLimit = raw.MemoryStats.Limit
	if raw.MemoryStats.Limit > 0 {
		m.MemPercent = float64(usage) / float64(raw.MemoryStats.Limit) * 100.0
	}

	// Network totals across all interfaces.
	for _, n := range raw.Networks {
		m.NetRx += n.RxBytes
		m.NetTx += n.TxBytes
	}

	// Block I/O totals.
	for _, b := range raw.BlkioStats.IoServiceBytesRecursive {
		switch strings.ToLower(b.Op) {
		case "read":
			m.BlkRead += b.Value
		case "write":
			m.BlkWrite += b.Value
		}
	}

	return m
}
