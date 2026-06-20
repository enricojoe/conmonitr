package docker

import (
	"bytes"
	"context"
	"io"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/pkg/stdcopy"
)

// LogLine is a single line of container output tagged with its source stream.
type LogLine struct {
	Stream    string `json:"stream"` // "stdout" or "stderr"
	Line      string `json:"line"`
	Timestamp int64  `json:"timestamp"` // unix milliseconds
}

// lineWriter buffers bytes and emits a LogLine to out for each complete line.
type lineWriter struct {
	stream string
	out    chan<- LogLine
	ctx    context.Context
	buf    []byte
}

func (w *lineWriter) Write(p []byte) (int, error) {
	w.buf = append(w.buf, p...)
	for {
		i := bytes.IndexByte(w.buf, '\n')
		if i < 0 {
			break
		}
		line := string(bytes.TrimRight(w.buf[:i], "\r"))
		w.buf = w.buf[i+1:]
		select {
		case w.out <- LogLine{Stream: w.stream, Line: line, Timestamp: time.Now().UnixMilli()}:
		case <-w.ctx.Done():
			return len(p), w.ctx.Err()
		}
	}
	return len(p), nil
}

// StreamLogs follows a container's logs, demultiplexing stdout/stderr for
// non-TTY containers, and sends each line to out. Blocks until ctx is
// cancelled or the stream ends.
func (s *Service) StreamLogs(ctx context.Context, id string, tail string, out chan<- LogLine) error {
	info, err := s.cli.ContainerInspect(ctx, id)
	if err != nil {
		return err
	}
	tty := info.Config != nil && info.Config.Tty

	reader, err := s.cli.ContainerLogs(ctx, id, container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
		Tail:       tail,
	})
	if err != nil {
		return err
	}
	defer reader.Close()

	stdout := &lineWriter{stream: "stdout", out: out, ctx: ctx}
	if tty {
		// TTY logs are not multiplexed; copy raw into the stdout line writer.
		_, err = io.Copy(stdout, reader)
	} else {
		stderr := &lineWriter{stream: "stderr", out: out, ctx: ctx}
		_, err = stdcopy.StdCopy(stdout, stderr, reader)
	}
	if err == io.EOF || ctx.Err() != nil {
		return nil
	}
	return err
}
