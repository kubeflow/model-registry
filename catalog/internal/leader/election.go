// Package leader provides distributed leader election with PostgreSQL-based locking.
package leader

import (
	"context"
	"errors"
	"fmt"
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
	onBecomeLeader func(context.Context)
	client         *pglock.Client
}

// NewLeaderElector creates a new LeaderElector instance.
//
// Parameters:
//   - gormDB: PostgreSQL database connection
//   - ctx: Base context for leadership contexts
//   - lockName: Unique identifier for the distributed lock
//   - lockDuration: How long the lock is held before expiring
//   - heartbeatFreq: How often to renew the lock while leader
//   - onBecomeLeader: Callback invoked when this instance becomes leader.
//     The function receives a context that is canceled when leadership is lost.
//     The callback should block until the context is canceled for graceful shutdown.
//
// Returns a configured LeaderElector, or an error if client creation fails.
func NewLeaderElector(
	gormDB *gorm.DB,
	ctx context.Context,
	lockName string,
	lockDuration time.Duration,
	heartbeatFreq time.Duration,
	onBecomeLeader func(context.Context),
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

	return &LeaderElector{
		ctx:            ctx,
		lockName:       lockName,
		lockDuration:   lockDuration,
		heartbeatFreq:  heartbeatFreq,
		onBecomeLeader: onBecomeLeader,
		client:         client,
	}, nil
}

// Run starts the leader election process with automatic retry.
//
// Continuously attempts to acquire leadership until the context cancels.
// Handles retry logic internally with exponential backoff.
//
// Behavior:
//  1. Acquires the distributed lock
//  2. On success, invokes onBecomeLeader callback with a leadership context
//  3. Renews the lock at heartbeatFreq intervals
//  4. On loss or error, retries with exponential backoff
//  5. Returns only when context cancels (graceful shutdown)
//
// Blocks until the context cancels. Returns context.Canceled on graceful shutdown.
//
// Thread-safety: This method is safe to call from multiple goroutines,
// but only one instance will acquire leadership at a time.
func (e *LeaderElector) Run(ctx context.Context) error {
	backoff := newBackoff()

	for {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		glog.Info("Attempting to acquire leadership...")
		err := e.runOnce(ctx)

		if errors.Is(err, context.Canceled) {
			glog.Info("Leader election canceled, shutting down")
			return err
		}

		if err != nil {
			delay := backoff.next()
			glog.Errorf("Leader election error: %v (retrying in %v)", err, delay)

			select {
			case <-time.After(delay):
				continue
			case <-ctx.Done():
				return ctx.Err()
			}
		}

		glog.Info("Leadership ended gracefully, attempting to reacquire")
		backoff.reset()
	}
}

// runOnce attempts to acquire leadership once and run the leader callback.
// Returns when context is canceled, lock is lost, or callback exits.
func (e *LeaderElector) runOnce(ctx context.Context) error {
	lock, err := e.client.AcquireContext(ctx, e.lockName)
	if err != nil {
		return fmt.Errorf("failed to acquire lock %q: %w", e.lockName, err)
	}

	glog.Infof("Successfully acquired leadership lock: %s", e.lockName)

	// Create a context that will be canceled when we lose leadership
	leaderCtx, cancelLeader := context.WithCancel(e.ctx)
	defer cancelLeader()

	// Channel to signal when onBecomeLeader returns
	leaderDone := make(chan struct{})

	// Start the leader callback in a goroutine
	go func() {
		defer close(leaderDone)
		defer func() {
			if r := recover(); r != nil {
				glog.Errorf("Leader callback panicked: %v", r)
			}
		}()
		e.onBecomeLeader(leaderCtx)
	}()

	ticker := time.NewTicker(e.heartbeatFreq)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			// Parent context canceled - graceful shutdown
			glog.Infof("Context canceled, releasing leadership lock: %s", e.lockName)
			cancelLeader() // Signal leader callback to stop
			<-leaderDone   // Wait for callback to finish

			if err := e.client.Release(lock); err != nil {
				glog.Errorf("Error releasing lock: %v", err)
			}
			return ctx.Err()

		case <-ticker.C:
			// Verify we still own the lock
			if err := e.client.SendHeartbeat(ctx, lock); err != nil {
				glog.Errorf("Lost leadership lock: %v", err)
				cancelLeader() // Signal leader callback to stop
				<-leaderDone   // Wait for callback to finish
				return fmt.Errorf("lost leadership: %w", err)
			}

		case <-leaderDone:
			// Leader callback exited unexpectedly
			glog.Infof("Leader callback exited, releasing lock: %s", e.lockName)
			if err := e.client.Release(lock); err != nil {
				glog.Errorf("Error releasing lock: %v", err)
			}
			return fmt.Errorf("leader callback exited unexpectedly")
		}
	}
}
