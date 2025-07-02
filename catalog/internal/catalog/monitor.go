package catalog

import (
	"fmt"
	"hash/crc32"
	"io"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"

	"github.com/fsnotify/fsnotify"
	"github.com/golang/glog"
)

// monitor sends events when the contents of a file have changed.
//
// Unfortunately, simply watching the file misses events for our primary case
// of k8s mounted configmaps because the files we're watching are actually
// symlinks which aren't modified:
//
//	drwxrwxrwx    1 root     root           138 Jul  2 15:45 .
//	drwxr-xr-x    1 root     root           116 Jul  2 15:52 ..
//	drwxr-xr-x    1 root     root            62 Jul  2 15:45 ..2025_07_02_15_45_09.2837733502
//	lrwxrwxrwx    1 root     root            32 Jul  2 15:45 ..data -> ..2025_07_02_15_45_09.2837733502
//	lrwxrwxrwx    1 root     root            26 Jul  2 13:18 sample-catalog.yaml -> ..data/sample-catalog.yaml
//	lrwxrwxrwx    1 root     root            19 Jul  2 13:18 sources.yaml -> ..data/sources.yaml
//
// Updates are written to a new directory and the ..data symlink is updated. No
// fsnotify events will ever be triggered for the YAML files.
//
// The approach taken here is to watch the directory containing the file for
// any change and then hash the contents of the file to avoid false-positives.
type monitor struct {
	watcher *fsnotify.Watcher
	closed  <-chan struct{}

	recordsMu sync.RWMutex
	records   map[string]map[string]*monitorRecord
}

var _monitor *monitor
var initMonitor sync.Once

// getMonitor returns a singleton monitor instance. Panics on failure.
func getMonitor() *monitor {
	initMonitor.Do(func() {
		var err error
		_monitor, err = newMonitor()
		if err != nil {
			panic(fmt.Sprintf("Unable to create file monitor: %v", err))
		}
	})
	if _monitor == nil {
		// Panic in case someone traps the panic that occurred during
		// initialization and tries to call this again.
		panic("Unable to get file monitor")
	}

	return _monitor
}

func newMonitor() (*monitor, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	m := &monitor{
		watcher: watcher,
		records: map[string]map[string]*monitorRecord{},
	}

	go m.monitor()
	return m, nil
}

// Close stops the monitor and waits for the background goroutine to exit.
//
// All channels returned by Path() will be closed.
func (m *monitor) Close() {
	select {
	case <-m.closed:
		// Already closed, nothing to do.
		return
	default:
		// Fallthrough
	}

	m.watcher.Close()
	<-m.closed

	m.recordsMu.Lock()
	defer m.recordsMu.Unlock()

	uniqCh := make(map[chan<- struct{}]struct{})
	for dir := range m.records {
		for file := range m.records[dir] {
			record, ok := m.records[dir][file]
			if !ok {
				continue
			}
			for _, ch := range record.channels {
				uniqCh[ch] = struct{}{}
			}
		}
	}
	for ch := range uniqCh {
		close(ch)
	}
	m.records = nil
}

// Path returns a channel that receives an event when the contents of a file
// change. The file does not need to exist before calling this method, however
// the provided path should only be a file or a symlink (not a directory,
// device, etc.). The returned channel will be closed when the monitor is
// closed.
func (m *monitor) Path(p string) (<-chan struct{}, error) {
	absPath, err := filepath.Abs(p)
	if err != nil {
		return nil, fmt.Errorf("abs: %w", err)
	}

	m.recordsMu.Lock()
	defer m.recordsMu.Unlock()

	dir, base := filepath.Split(absPath)
	dir = filepath.Clean(dir)

	err = m.watcher.Add(dir)
	if err != nil {
		return nil, fmt.Errorf("unable to watch directory %q: %w", dir, err)
	}

	if _, exists := m.records[dir]; !exists {
		m.records[dir] = make(map[string]*monitorRecord, 1)
	}

	ch := make(chan struct{}, 1)

	if _, exists := m.records[dir][base]; !exists {
		m.records[dir][base] = &monitorRecord{
			channels: []chan<- struct{}{ch},
		}
	} else {
		r := m.records[dir][base]
		r.channels = append(r.channels, ch)
	}
	m.records[dir][base].updateHash(filepath.Join(dir, base))

	return ch, nil
}

func (m *monitor) monitor() {
	closed := make(chan struct{})
	m.closed = closed
	defer close(closed)

	for {
		select {
		case err, ok := <-m.watcher.Errors:
			if !ok {
				return
			}

			glog.Errorf("fsnotify error: %v", err)
		case e, ok := <-m.watcher.Events:
			if !ok {
				return
			}

			glog.V(2).Infof("fsnotify.Event: %v", e)

			switch e.Op {
			case fsnotify.Create, fsnotify.Write:
				// Fallthrough
			default:
				// Ignore fsnotify.Remove, fsnotify.Rename and fsnotify.Chmod
				continue
			}

			func() {
				m.recordsMu.RLock()
				defer m.recordsMu.RUnlock()

				dir := filepath.Dir(e.Name)

				dc := m.records[dir]
				if dc == nil {
					return
				}

				for base, record := range dc {
					path := filepath.Join(dir, base)
					if !record.updateHash(path) {
						continue
					}
					for _, ch := range record.channels {
						// Send the event, ignore any that would block.
						select {
						case ch <- struct{}{}:
						default:
							glog.Errorf("monitor: missed event for path %s", path)
						}
					}
				}
			}()
		}
	}
}

type monitorRecord struct {
	channels []chan<- struct{}
	hash     uint32
}

// updateHash recalculates the hash and returns true if it has changed.
func (mr *monitorRecord) updateHash(path string) bool {
	newHash := mr.calculateHash(path)
	oldHash := atomic.SwapUint32(&mr.hash, newHash)
	return oldHash != newHash
}

func (monitorRecord) calculateHash(path string) uint32 {
	fh, err := os.Open(path)
	if err != nil {
		return 0
	}
	defer fh.Close()

	h := crc32.NewIEEE()
	_, err = io.Copy(h, fh)
	if err != nil {
		return 0
	}
	return h.Sum32()
}
