package log

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"sync"

	eventv1 "github.com/sparrow-community/sparrow/protocol/gen/go/event/v1"
	"google.golang.org/protobuf/encoding/protodelim"
	"google.golang.org/protobuf/proto"
)

// File is an append-only EventLog backed by length-delimited protobuf records.
type File struct {
	mu   sync.Mutex
	path string
	f    *os.File
}

type countingReader struct {
	r io.Reader
	n int64
}

func (c *countingReader) Read(p []byte) (int, error) {
	n, err := c.r.Read(p)
	c.n += int64(n)
	return n, err
}

// OpenFile opens (or creates) a file-backed event log at path.
func OpenFile(path string) (*File, error) {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0o644)
	if err != nil {
		return nil, fmt.Errorf("open event log: %w", err)
	}
	if _, err := f.Seek(0, io.SeekEnd); err != nil {
		_ = f.Close()
		return nil, fmt.Errorf("seek event log: %w", err)
	}
	return &File{path: path, f: f}, nil
}

func (fl *File) Close() error {
	fl.mu.Lock()
	defer fl.mu.Unlock()
	if fl.f == nil {
		return nil
	}
	err := fl.f.Close()
	fl.f = nil
	return err
}

func (fl *File) Append(_ context.Context, e *eventv1.Event) (int64, error) {
	if e == nil {
		return 0, fmt.Errorf("event is nil")
	}
	clone, ok := proto.Clone(e).(*eventv1.Event)
	if !ok || clone == nil {
		return 0, fmt.Errorf("clone event")
	}

	fl.mu.Lock()
	defer fl.mu.Unlock()
	if fl.f == nil {
		return 0, fmt.Errorf("event log closed")
	}

	if _, err := protodelim.MarshalTo(fl.f, clone); err != nil {
		return 0, fmt.Errorf("append event: %w", err)
	}
	if err := fl.f.Sync(); err != nil {
		return 0, fmt.Errorf("sync event log: %w", err)
	}
	pos, err := fl.f.Seek(0, io.SeekCurrent)
	if err != nil {
		return 0, err
	}
	return pos, nil
}

func (fl *File) ReadByInstance(ctx context.Context, processInstanceID string) ([]*eventv1.Event, error) {
	all, err := fl.ReadAll(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]*eventv1.Event, 0)
	for _, e := range all {
		if e.GetProcessInstanceId() == processInstanceID {
			out = append(out, e)
		}
	}
	return out, nil
}

func (fl *File) ReadAll(_ context.Context) ([]*eventv1.Event, error) {
	fl.mu.Lock()
	defer fl.mu.Unlock()
	if fl.f == nil {
		return nil, fmt.Errorf("event log closed")
	}

	if _, err := fl.f.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}

	cr := &countingReader{r: fl.f}
	r := bufio.NewReader(cr)
	out := make([]*eventv1.Event, 0, 64)
	var lastGood int64
	for {
		ev := &eventv1.Event{}
		if err := protodelim.UnmarshalFrom(r, ev); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			if errors.Is(err, io.ErrUnexpectedEOF) {
				// Crash mid-append: drop the incomplete tail so the next Append is valid.
				if err := fl.f.Truncate(lastGood); err != nil {
					return nil, fmt.Errorf("truncate incomplete event: %w", err)
				}
				break
			}
			return nil, fmt.Errorf("read event: %w", err)
		}
		lastGood = cr.n - int64(r.Buffered())
		out = append(out, ev)
	}
	if _, err := fl.f.Seek(0, io.SeekEnd); err != nil {
		return nil, err
	}
	return out, nil
}
