package catalog

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var enableMonitorTests = flag.Bool("monitor-tests", false, "Enable flaky monitor tests (disabled by default)")

// skipFlaky checks if flaky monitor tests should be skipped.
// These tests are disabled by default due to timing sensitivity in CI environments.
// Run with -monitor-tests flag to enable them.
func skipFlaky(t *testing.T) {
	t.Helper()
	if !*enableMonitorTests {
		t.Skip("Skipping flaky monitor test. Run with -monitor-tests flag to enable.")
	}
}

func TestMonitor(t *testing.T) {
	skipFlaky(t)
	assert := assert.New(t)

	mon, err := newMonitor()
	if !assert.NoError(err) {
		return
	}
	mon.Pause = 0

	tmpDir := t.TempDir()
	fileA := filepath.Join(tmpDir, "a")
	fileB := filepath.Join(tmpDir, "b")
	fileC := filepath.Join(tmpDir, "c")

	_watchMonitor := func(ch <-chan struct{}, err error) *monitorWatcher {
		if err != nil {
			t.Fatalf("watchMonitor passed error %v", err)
		}
		return watchMonitor(ch)
	}

	a := _watchMonitor(mon.Path(t.Context(), fileA))
	b := _watchMonitor(mon.Path(t.Context(), fileB))

	updateFile(t, fileA)
	a.AssertCount(t, 1)
	b.AssertCount(t, 0, "unchanged file should not have any events")

	a.Reset()
	updateFile(t, fileB)
	b.AssertCount(t, 1)
	updateFile(t, fileB)
	b.AssertCount(t, 2)
	a.AssertCount(t, 0, "unchanged file should not have any events")

	b.Reset()
	updateFile(t, fileC)
	a.AssertCount(t, 0, "unchanged file should not have an event")
	b.AssertCount(t, 0, "unchanged file should not have an event")

	// Ensure that Close doesn't hang.
	finished := make(chan struct{})
	go func() {
		defer close(finished)
		mon.Close()
	}()
	assert.Eventually(func() bool {
		select {
		case <-finished:
			return true
		default:
			return false
		}
	}, time.Second, 50*time.Millisecond)

	// Verify that the monitor channels closed.
	assert.True(a.Done())
	assert.True(b.Done())
}

func TestMonitorSymlinks(t *testing.T) {
	skipFlaky(t)
	assert := assert.New(t)

	tmpDir := t.TempDir()
	mon, err := newMonitor()
	if !assert.NoError(err) {
		return
	}
	mon.Pause = 0
	defer mon.Close()

	// Watch the files on the published path.
	_watchMonitor := func(ch <-chan struct{}, err error) *monitorWatcher {
		if err != nil {
			t.Fatalf("watchMonitor passed error %v", err)
		}
		return watchMonitor(ch)
	}

	a := _watchMonitor(mon.Path(t.Context(), filepath.Join(tmpDir, "a")))
	b := _watchMonitor(mon.Path(t.Context(), filepath.Join(tmpDir, "b")))

	// Set up a directory structure with symlinks like k8s does for mounted
	// configmaps.
	// a -> latest/a, b -> latest/b, latest -> v1
	assert.NoError(os.Mkdir(filepath.Join(tmpDir, "v1"), 0777))
	updateFile(t, filepath.Join(tmpDir, "v1", "a"), "foo")
	updateFile(t, filepath.Join(tmpDir, "v1", "b"), "bar")
	assert.NoError(os.Symlink("v1", filepath.Join(tmpDir, "latest")))
	assert.NoError(os.Symlink(filepath.Join("latest", "a"), filepath.Join(tmpDir, "a")))
	assert.NoError(os.Symlink(filepath.Join("latest", "b"), filepath.Join(tmpDir, "b")))

	a.AssertCount(t, 1)
	b.AssertCount(t, 1)
	a.Reset()
	b.Reset()

	// Make a new version directory
	os.Mkdir(filepath.Join(tmpDir, "v2"), 0777)
	updateFile(t, filepath.Join(tmpDir, "v2", "a"), "UPDATED")
	updateFile(t, filepath.Join(tmpDir, "v2", "b"), "bar")

	a.AssertCount(t, 0)
	b.AssertCount(t, 0)
	a.Reset()
	b.Reset()

	// Update the symlink to point to the new version:
	assert.NoError(os.Rename(filepath.Join(tmpDir, "latest"), filepath.Join(tmpDir, "latest_tmp")))
	assert.NoError(os.Symlink(filepath.Join("v2"), filepath.Join(tmpDir, "latest")))
	assert.NoError(os.Remove(filepath.Join(tmpDir, "latest_tmp")))
	assert.NoError(os.RemoveAll(filepath.Join(tmpDir, "v1")))

	a.AssertCount(t, 1)
	b.AssertCount(t, 0)
}

type monitorWatcher struct {
	count int32
	done  int32
}

func (mw *monitorWatcher) Reset() {
	atomic.StoreInt32(&mw.count, 0)
}

func (mw *monitorWatcher) AssertCount(t *testing.T, expected int, args ...any) bool {
	t.Helper()
	// Use 5 second timeout to handle slow CI runners where file system events
	// can take longer to propagate. The test still returns immediately when
	// the condition is met, so this doesn't slow down fast runs.
	return assert.Eventually(t, func() bool {
		return int(atomic.LoadInt32(&mw.count)) == expected
	}, 5*time.Second, 10*time.Millisecond, args...)
}

func (mw *monitorWatcher) Count() int {
	return int(atomic.LoadInt32(&mw.count))
}

func (mw *monitorWatcher) Done() bool {
	return atomic.LoadInt32(&mw.done) != 0
}

func watchMonitor(ch <-chan struct{}) *monitorWatcher {
	mw := &monitorWatcher{}

	go func() {
		defer atomic.StoreInt32(&mw.done, 1)
		for range ch {
			atomic.AddInt32(&mw.count, 1)
		}
	}()

	return mw
}

func updateFile(t *testing.T, path string, contents ...string) {
	fh, err := os.Create(path)
	if err != nil {
		t.Fatalf("unable to open %q: %v", path, err)
	}
	if len(contents) == 0 {
		fmt.Fprintf(fh, "%s\n", time.Now())
	} else {
		for _, line := range contents {
			fmt.Fprintf(fh, "%s\n", line)
		}
	}
	fh.Close()
}
