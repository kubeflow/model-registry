package cmd

import (
	"os"
	"testing"
	"time"
)

func TestGetLeaderElectionConfig(t *testing.T) {
	tests := []struct {
		name                 string
		lockDurationEnv      string
		heartbeatEnv         string
		expectedLockDuration time.Duration
		expectedHeartbeat    time.Duration
	}{
		{
			name:                 "defaults when no env vars set",
			lockDurationEnv:      "",
			heartbeatEnv:         "",
			expectedLockDuration: defaultLeaderLockDuration,
			expectedHeartbeat:    defaultLeaderHeartbeat,
		},
		{
			name:                 "custom valid values",
			lockDurationEnv:      "30s",
			heartbeatEnv:         "10s",
			expectedLockDuration: 30 * time.Second,
			expectedHeartbeat:    10 * time.Second,
		},
		{
			name:                 "fast failover for local dev",
			lockDurationEnv:      "10s",
			heartbeatEnv:         "3s",
			expectedLockDuration: 10 * time.Second,
			expectedHeartbeat:    3 * time.Second,
		},
		{
			name:                 "long lease for production",
			lockDurationEnv:      "120s",
			heartbeatEnv:         "30s",
			expectedLockDuration: 120 * time.Second,
			expectedHeartbeat:    30 * time.Second,
		},
		{
			name:                 "invalid lock duration uses default lock, keeps valid heartbeat",
			lockDurationEnv:      "invalid",
			heartbeatEnv:         "10s",
			expectedLockDuration: defaultLeaderLockDuration, // falls back to default
			expectedHeartbeat:    10 * time.Second,          // uses parsed value
		},
		{
			name:                 "invalid heartbeat uses default heartbeat, keeps valid lock",
			lockDurationEnv:      "30s",
			heartbeatEnv:         "invalid",
			expectedLockDuration: 30 * time.Second,       // uses parsed value
			expectedHeartbeat:    defaultLeaderHeartbeat, // falls back to default
		},
		{
			name:                 "heartbeat exceeds lock duration/2 uses defaults",
			lockDurationEnv:      "30s",
			heartbeatEnv:         "20s", // 20s > 15s (half of 30s)
			expectedLockDuration: defaultLeaderLockDuration,
			expectedHeartbeat:    defaultLeaderHeartbeat,
		},
		{
			name:                 "heartbeat exactly at lock duration/2 is valid",
			lockDurationEnv:      "30s",
			heartbeatEnv:         "15s",
			expectedLockDuration: 30 * time.Second,
			expectedHeartbeat:    15 * time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear environment
			os.Unsetenv(envLeaderLockDuration)
			os.Unsetenv(envLeaderHeartbeat)

			// Set test environment
			if tt.lockDurationEnv != "" {
				os.Setenv(envLeaderLockDuration, tt.lockDurationEnv)
				defer os.Unsetenv(envLeaderLockDuration)
			}
			if tt.heartbeatEnv != "" {
				os.Setenv(envLeaderHeartbeat, tt.heartbeatEnv)
				defer os.Unsetenv(envLeaderHeartbeat)
			}

			// Get configuration
			lockDuration, heartbeat := getLeaderElectionConfig()

			// Verify
			if lockDuration != tt.expectedLockDuration {
				t.Errorf("lock duration = %v, want %v", lockDuration, tt.expectedLockDuration)
			}
			if heartbeat != tt.expectedHeartbeat {
				t.Errorf("heartbeat = %v, want %v", heartbeat, tt.expectedHeartbeat)
			}

			// Verify pglock requirement: heartbeat <= lockDuration/2
			if heartbeat > lockDuration/2 {
				t.Errorf("heartbeat (%v) exceeds half of lock duration (%v), violates pglock requirement", heartbeat, lockDuration)
			}
		})
	}
}
