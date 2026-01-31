package leader_test

import (
	"context"
	"os"
	"sync/atomic"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/kubeflow/model-registry/catalog/internal/db/service"
	"github.com/kubeflow/model-registry/catalog/internal/leader"
	"github.com/kubeflow/model-registry/internal/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	os.Exit(testutils.TestMainPostgresHelper(m))
}

func TestLeaderElector_Run(t *testing.T) {
	db, cleanup := testutils.SetupPostgresWithMigrations(t, service.DatastoreSpec())
	defer cleanup()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var becameLeader atomic.Bool
	var lostLeadership atomic.Bool

	onBecomeLeader := func(ctx context.Context) {
		becameLeader.Store(true)
		<-ctx.Done() // Block until leadership is lost
		lostLeadership.Store(true)
	}

	elector, err := leader.NewLeaderElector(
		db,
		ctx,
		"test-lock",
		10*time.Second,
		1*time.Second,
		onBecomeLeader,
	)
	require.NoError(t, err)

	errCh := make(chan error, 1)
	go func() {
		err := elector.Run(ctx)
		if err != nil {
			t.Logf("Run returned error: %v", err)
		}
		errCh <- err
	}()

	// Wait for leadership acquisition
	require.Eventually(t, func() bool {
		select {
		case err := <-errCh:
			t.Fatalf("Run exited early with error: %v", err)
		default:
		}
		return becameLeader.Load()
	}, 3*time.Second, 100*time.Millisecond, "Should acquire leadership")

	// Cancel context to trigger leadership loss
	cancel()

	// Wait for graceful shutdown
	select {
	case err := <-errCh:
		assert.ErrorIs(t, err, context.Canceled)
	case <-time.After(2 * time.Second):
		t.Fatal("Timeout waiting for Run to complete")
	}

	assert.True(t, lostLeadership.Load(), "Should have lost leadership gracefully")
}

// TestMultiPodLeaderElection verifies that only one pod can hold the leader lock at a time.
// This simulates a multi-pod deployment scenario where:
// - Multiple pods attempt to acquire leadership
// - Only one pod becomes leader at a time
// - When the leader releases the lock, another pod can acquire it
func TestMultiPodLeaderElection(t *testing.T) {
	db, cleanup := testutils.SetupPostgresWithMigrations(t, service.DatastoreSpec())
	defer cleanup()

	// Use a longer timeout to account for container startup and leader transitions
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Track which "pod" is currently the leader
	var currentLeader atomic.Int32
	var pod1BecameLeader, pod2BecameLeader atomic.Bool

	// Pod 1 callback
	onBecomeLeaderPod1 := func(ctx context.Context) {
		currentLeader.Store(1)
		pod1BecameLeader.Store(true)
		t.Log("Pod 1 became leader")
		<-ctx.Done() // Block until leadership is lost
		t.Log("Pod 1 lost leadership")
	}

	// Pod 2 callback
	onBecomeLeaderPod2 := func(ctx context.Context) {
		currentLeader.Store(2)
		pod2BecameLeader.Store(true)
		t.Log("Pod 2 became leader")
		<-ctx.Done() // Block until leadership is lost
		t.Log("Pod 2 lost leadership")
	}

	// Create two leader electors (simulating two pods)
	elector1, err := leader.NewLeaderElector(
		db,
		ctx,
		"catalog-leader",
		5*time.Second, // Lock duration
		1*time.Second, // Heartbeat frequency
		onBecomeLeaderPod1,
	)
	require.NoError(t, err)

	elector2, err := leader.NewLeaderElector(
		db,
		ctx,
		"catalog-leader", // Same lock name
		5*time.Second,
		1*time.Second,
		onBecomeLeaderPod2,
	)
	require.NoError(t, err)

	// Start pod 1 attempting to acquire leadership
	errCh1 := make(chan error, 1)
	ctx1, cancel1 := context.WithCancel(ctx)
	defer cancel1()
	go func() {
		errCh1 <- elector1.Run(ctx1)
	}()

	// Wait for pod 1 to become leader
	require.Eventually(t, func() bool {
		select {
		case err := <-errCh1:
			t.Fatalf("Pod 1 Run exited early with error: %v", err)
		default:
		}
		return pod1BecameLeader.Load()
	}, 3*time.Second, 100*time.Millisecond, "Pod 1 should acquire leadership")

	assert.Equal(t, int32(1), currentLeader.Load(), "Pod 1 should be the leader")

	// Start pod 2 attempting to acquire leadership (should block since pod 1 holds the lock)
	errCh2 := make(chan error, 1)
	ctx2, cancel2 := context.WithCancel(ctx)
	defer cancel2()
	go func() {
		errCh2 <- elector2.Run(ctx2)
	}()

	// Give pod 2 some time to attempt acquisition - it should not become leader
	time.Sleep(2 * time.Second)
	assert.False(t, pod2BecameLeader.Load(), "Pod 2 should not acquire leadership while pod 1 holds the lock")
	assert.Equal(t, int32(1), currentLeader.Load(), "Pod 1 should still be the leader")

	// Release pod 1's leadership
	t.Log("Releasing pod 1's leadership")
	cancel1()

	// Wait for pod 1's Run to complete
	select {
	case err := <-errCh1:
		assert.ErrorIs(t, err, context.Canceled, "Pod 1 should exit with context.Canceled")
	case <-time.After(3 * time.Second):
		t.Fatal("Timeout waiting for pod 1 to release leadership")
	}

	// Now pod 2 should be able to acquire leadership
	t.Log("Waiting for pod 2 to acquire leadership...")
	require.Eventually(t, func() bool {
		select {
		case err := <-errCh2:
			t.Fatalf("Pod 2 Run exited early with error: %v", err)
		default:
		}
		acquired := pod2BecameLeader.Load()
		if !acquired {
			t.Log("Pod 2 has not yet acquired leadership, waiting...")
		}
		return acquired
	}, 10*time.Second, 200*time.Millisecond, "Pod 2 should acquire leadership after pod 1 releases")

	assert.Equal(t, int32(2), currentLeader.Load(), "Pod 2 should be the leader")

	// Clean up pod 2
	cancel2()
	select {
	case err := <-errCh2:
		assert.ErrorIs(t, err, context.Canceled, "Pod 2 should exit with context.Canceled")
	case <-time.After(3 * time.Second):
		t.Fatal("Timeout waiting for pod 2 to release leadership")
	}
}
