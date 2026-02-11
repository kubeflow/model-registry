// Package leader provides distributed leader election with PostgreSQL-based locking.
package leader

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"cirello.io/pglock"
	"github.com/golang/glog"
	"gorm.io/gorm"
)

// backoff provides exponential backoff for retry logic.
type backoff struct {
	current time.Duration
	max     time.Duration
}

// newBackoff creates a backoff with 1s initial delay and 30s maximum.
func newBackoff() *backoff {
	return &backoff{
		current: 1 * time.Second,
		max:     30 * time.Second,
	}
}

// next returns the current delay, then doubles it (up to max).
func (b *backoff) next() time.Duration {
	delay := b.current
	b.current = min(b.current*2, b.max)
	return delay
}

func (b *backoff) reset() {
	b.current = 1 * time.Second
}

// LeaderElector manages leader election for distributed services
// using pglock to coordinate leadership across multiple instances.
type LeaderElector struct {
	ctx            context.Context
	lockName       string
	lockDuration   time.Duration
	heartbeatFreq  time.Duration
	onBecomeLeader []func(context.Context)
	client         *pglock.Client
	mu             sync.Mutex

	// Leadership state (protected by mu)
	isLeader     bool
	leaderCtx    context.Context
	cancelLeader func()

	// Background goroutine tracking
	done chan struct{}
	err  error

	// Dynamic callback tracking
	activeCallbacks sync.WaitGroup
}

// NewLeaderElector creates a new LeaderElector instance.
//
// Parameters:
//   - gormDB: PostgreSQL database connection
//   - ctx: Base context for leadership contexts
//   - lockName: Unique identifier for the distributed lock
//   - lockDuration: How long the lock is held before expiring
//   - heartbeatFreq: How often to renew the lock while leader
//
// Returns a configured LeaderElector, or an error if client creation fails.
// Use OnBecomeLeader to register callbacks that will be invoked when leadership is acquired.
func NewLeaderElector(
	gormDB *gorm.DB,
	ctx context.Context,
	lockName string,
	lockDuration time.Duration,
	heartbeatFreq time.Duration,
) (*LeaderElector, error) {
	if gormDB.Name() != "postgres" {
		return nil, errors.New("not a postgres database handle")
	}

	db, err := gormDB.DB()
	if err != nil {
		return nil, fmt.Errorf("unable to get sql.DB from GORM: %w", err)
	}

	client, err := pglock.UnsafeNew(
		db,
		pglock.WithLeaseDuration(lockDuration),
		pglock.WithHeartbeatFrequency(heartbeatFreq),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create pglock client: %w", err)
	}

	err = client.TryCreateTable()
	if err != nil {
		return nil, err
	}

	e := &LeaderElector{
		ctx:           ctx,
		lockName:      lockName,
		lockDuration:  lockDuration,
		heartbeatFreq: heartbeatFreq,
		client:        client,
		done:          make(chan struct{}),
	}

	// Start the background goroutine immediately
	go e.run()

	return e, nil
}

// startCallback launches a callback in a goroutine with panic recovery and proper cleanup.
func (e *LeaderElector) startCallback(leaderCtx context.Context, idx int, cb func(context.Context)) {
	e.activeCallbacks.Add(1)
	go func() {
		defer e.activeCallbacks.Done()
		defer func() {
			if r := recover(); r != nil {
				glog.Errorf("Leader callback %d panicked: %v", idx, r)
			}
		}()
		cb(leaderCtx)
		// Only warn if context wasn't cancelled
		if leaderCtx.Err() == nil {
			glog.Warningf("Leader callback %d exited early", idx)
		}
	}()
}

// OnBecomeLeader registers a callback to be invoked when this instance becomes leader.
// The callback receives a context that is canceled when leadership is lost.
// The callback should block until the context is canceled for graceful shutdown.
// Multiple callbacks can be registered and will be executed concurrently.
//
// If this instance is already the leader when this method is called, the callback
// will be invoked immediately with the current leadership context.
func (e *LeaderElector) OnBecomeLeader(callback func(context.Context)) {
	e.mu.Lock()
	e.onBecomeLeader = append(e.onBecomeLeader, callback)

	// If already leader, start callback immediately
	if e.isLeader {
		leaderCtx := e.leaderCtx
		idx := len(e.onBecomeLeader) - 1
		e.mu.Unlock()
		e.startCallback(leaderCtx, idx, callback)
		return
	}
	e.mu.Unlock()
}

// Wait blocks until the background goroutine exits and returns any error.
// This replaces the old Run() method in the new API.
func (e *LeaderElector) Wait() error {
	<-e.done
	return e.err
}

// run starts the leader election process with automatic retry.
// This runs in a background goroutine started by NewLeaderElector.
//
// Continuously attempts to acquire leadership until the context cancels.
// Handles retry logic internally with exponential backoff.
//
// Behavior:
//  1. Acquires the distributed lock
//  2. On success, invokes all registered callbacks concurrently with a leadership context
//  3. Renews the lock at heartbeatFreq intervals
//  4. On loss or error, retries with exponential backoff
//  5. Returns only when context cancels (graceful shutdown)
func (e *LeaderElector) run() {
	defer close(e.done)

	ctx := e.ctx
	backoff := newBackoff()

	for {
		if ctx.Err() != nil {
			e.err = ctx.Err()
			return
		}

		glog.Info("Attempting to acquire leadership...")
		err := e.runOnce(ctx)

		if errors.Is(err, context.Canceled) {
			glog.Info("Leader election canceled, shutting down")
			e.err = err
			return
		}

		if err != nil {
			delay := backoff.next()
			glog.Errorf("Leader election error: %v (retrying in %v)", err, delay)

			select {
			case <-time.After(delay):
				continue
			case <-ctx.Done():
				e.err = ctx.Err()
				return
			}
		}

		glog.Info("Leadership ended gracefully, attempting to reacquire")
		backoff.reset()
	}
}

// runOnce attempts to acquire leadership once and run the leader callbacks.
// Returns when context is canceled or lock is lost.
// Per the plan, callbacks exiting early no longer causes lock release.
func (e *LeaderElector) runOnce(ctx context.Context) error {
	lock, err := e.client.AcquireContext(ctx, e.lockName)
	if err != nil {
		return fmt.Errorf("failed to acquire lock %q: %w", e.lockName, err)
	}

	glog.Infof("Successfully acquired leadership lock: %s", e.lockName)

	// Create a context that will be canceled when we lose leadership
	leaderCtx, cancelLeader := context.WithCancel(e.ctx)
	defer cancelLeader()

	// Set leadership state and store leaderCtx before starting callbacks
	e.mu.Lock()
	e.isLeader = true
	e.leaderCtx = leaderCtx
	e.cancelLeader = cancelLeader

	// Get snapshot of callbacks
	callbacks := make([]func(context.Context), len(e.onBecomeLeader))
	copy(callbacks, e.onBecomeLeader)
	e.mu.Unlock()

	// Start all callbacks concurrently using startCallback
	for i, callback := range callbacks {
		e.startCallback(leaderCtx, i, callback)
	}

	ticker := time.NewTicker(e.heartbeatFreq)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			// Parent context canceled - graceful shutdown
			glog.Infof("Context canceled, releasing leadership lock: %s", e.lockName)
			cancelLeader()           // Signal ALL callbacks to stop
			e.activeCallbacks.Wait() // Wait for ALL callbacks to finish

			// Clear leadership state
			e.mu.Lock()
			e.isLeader = false
			e.leaderCtx = nil
			e.cancelLeader = nil
			e.mu.Unlock()

			if err := e.client.Release(lock); err != nil {
				glog.Errorf("Error releasing lock: %v", err)
			}
			return ctx.Err()

		case <-ticker.C:
			// Verify we still own the lock
			if err := e.client.SendHeartbeat(ctx, lock); err != nil {
				glog.Errorf("Lost leadership lock: %v", err)
				cancelLeader()           // Signal ALL callbacks to stop
				e.activeCallbacks.Wait() // Wait for ALL callbacks to finish

				// Clear leadership state
				e.mu.Lock()
				e.isLeader = false
				e.leaderCtx = nil
				e.cancelLeader = nil
				e.mu.Unlock()

				return fmt.Errorf("lost leadership: %w", err)
			}
		}
	}
}
