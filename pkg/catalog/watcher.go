package catalog

import (
	"context"
	"os"
	"sync"
	"time"
)

// FileWatcher watches a file for changes using polling.
// This is a simple implementation that can be used if more sophisticated
// file watching (like fsnotify) is not available.
type FileWatcher struct {
	mu          sync.Mutex
	path        string
	lastModTime time.Time
	interval    time.Duration
}

// NewFileWatcher creates a new file watcher.
func NewFileWatcher(path string, interval time.Duration) *FileWatcher {
	if interval == 0 {
		interval = 5 * time.Second
	}

	w := &FileWatcher{
		path:     path,
		interval: interval,
	}

	if info, err := os.Stat(path); err == nil {
		w.lastModTime = info.ModTime()
	}

	return w
}

// Watch returns a channel that receives a value whenever the file changes.
// The channel is closed when the context is canceled.
func (w *FileWatcher) Watch(ctx context.Context) <-chan struct{} {
	ch := make(chan struct{})

	go func() {
		defer close(ch)

		ticker := time.NewTicker(w.interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if w.hasChanged() {
					select {
					case ch <- struct{}{}:
					case <-ctx.Done():
						return
					}
				}
			}
		}
	}()

	return ch
}

func (w *FileWatcher) hasChanged() bool {
	w.mu.Lock()
	defer w.mu.Unlock()

	info, err := os.Stat(w.path)
	if err != nil {
		return false
	}

	if info.ModTime().After(w.lastModTime) {
		w.lastModTime = info.ModTime()
		return true
	}

	return false
}
