package test

import (
	"context"
	"testing"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	natslib "github.com/pondigo/xk6-nats"
)

func TestConnectionLifecycle(t *testing.T) {
	// Test connection creation and lifecycle
	_ = &MockVU{}

	// Test connection options validation
	opts := natslib.ConnectionOptions{
		URLs:           []string{nats.DefaultURL},
		MaxReconnects:  5,
		ReconnectWait:  2,
		PingInterval:   60,
		MaxPingsOut:    2,
		AllowReconnect: true,
	}

	err := natslib.ValidateConnectionOptions(opts)
	assert.NoError(t, err)

	// Note: We can't test actual connection without a NATS server
	// This would require setting up a test server
}

func TestStreamConfigConversion(t *testing.T) {
	tests := []struct {
		name     string
		config   natslib.StreamConfig
		expected *nats.StreamConfig
	}{
		{
			name: "limits retention",
			config: natslib.StreamConfig{
				Name:      "TEST_STREAM",
				Subjects:  []string{"test.>"},
				Retention: "limits",
				Storage:   "file",
				Discard:   "old",
				Replicas:  1,
				MaxBytes:  1024 * 1024,
				MaxMsgs:   1000,
				MaxAge:    3600,
			},
			expected: &nats.StreamConfig{
				Name:      "TEST_STREAM",
				Subjects:  []string{"test.>"},
				Retention: nats.LimitsPolicy,
				Storage:   nats.FileStorage,
				Discard:   nats.DiscardOld,
				Replicas:  1,
				MaxBytes:  1024 * 1024,
				MaxMsgs:   1000,
				MaxAge:    3600 * time.Second,
			},
		},
		{
			name: "interest retention",
			config: natslib.StreamConfig{
				Name:      "TEST_STREAM",
				Subjects:  []string{"test.>"},
				Retention: "interest",
				Storage:   "memory",
				Discard:   "new",
				Replicas:  1,
			},
			expected: &nats.StreamConfig{
				Name:      "TEST_STREAM",
				Subjects:  []string{"test.>"},
				Retention: nats.InterestPolicy,
				Storage:   nats.MemoryStorage,
				Discard:   nats.DiscardNew,
				Replicas:  1,
			},
		},
		{
			name: "workqueue retention",
			config: natslib.StreamConfig{
				Name:      "TEST_STREAM",
				Subjects:  []string{"test.>"},
				Retention: "workqueue",
				Storage:   "memory",
				Replicas:  1,
			},
			expected: &nats.StreamConfig{
				Name:      "TEST_STREAM",
				Subjects:  []string{"test.>"},
				Retention: nats.WorkQueuePolicy,
				Storage:   nats.MemoryStorage,
				Replicas:  1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This would test the conversion logic in jetstream.go
			// For now, we just validate the config
			err := natslib.ValidateStreamConfig(tt.config)
			assert.NoError(t, err)
		})
	}
}

func TestConsumerConfigConversion(t *testing.T) {
	tests := []struct {
		name     string
		config   natslib.ConsumerConfig
		expected *nats.ConsumerConfig
	}{
		{
			name: "explicit ack policy",
			config: natslib.ConsumerConfig{
				Stream:        "TEST_STREAM",
				Durable:       "TEST_CONSUMER",
				DeliverPolicy: "all",
				AckPolicy:     "explicit",
				AckWait:       30,
				MaxDeliver:    3,
				ReplayPolicy:  "instant",
			},
			expected: &nats.ConsumerConfig{
				Durable:       "TEST_CONSUMER",
				DeliverPolicy: nats.DeliverAllPolicy,
				AckPolicy:     nats.AckExplicitPolicy,
				AckWait:       30 * time.Second,
				MaxDeliver:    3,
				ReplayPolicy:  nats.ReplayInstantPolicy,
			},
		},
		{
			name: "none ack policy",
			config: natslib.ConsumerConfig{
				Stream:        "TEST_STREAM",
				Durable:       "TEST_CONSUMER",
				DeliverPolicy: "last",
				AckPolicy:     "none",
				ReplayPolicy:  "original",
			},
			expected: &nats.ConsumerConfig{
				Durable:       "TEST_CONSUMER",
				DeliverPolicy: nats.DeliverLastPolicy,
				AckPolicy:     nats.AckNonePolicy,
				ReplayPolicy:  nats.ReplayOriginalPolicy,
			},
		},
		{
			name: "all ack policy",
			config: natslib.ConsumerConfig{
				Stream:        "TEST_STREAM",
				Durable:       "TEST_CONSUMER",
				DeliverPolicy: "new",
				AckPolicy:     "all",
			},
			expected: &nats.ConsumerConfig{
				Durable:       "TEST_CONSUMER",
				DeliverPolicy: nats.DeliverNewPolicy,
				AckPolicy:     nats.AckAllPolicy,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This would test the conversion logic in consumer.go
			// For now, we just validate the config
			err := natslib.ValidateConsumerConfig(tt.config)
			assert.NoError(t, err)
		})
	}
}

func TestMetricsCreation(t *testing.T) {
	vu := &MockVU{
		ctx: context.Background(),
	}

	// Test metrics creation
	metrics, err := natslib.NewNatsMetrics(vu)
	require.NoError(t, err)
	require.NotNil(t, metrics)

	// Test registry
	registry := metrics.Registry()
	assert.NotNil(t, registry)

	// Test placeholder methods (they should not panic)
	assert.NotPanics(t, func() {
		metrics.RecordConnectionEstablished()
		metrics.RecordConnectionClosed()
		metrics.RecordConnectionFailed()
		metrics.RecordReconnect()
		metrics.RecordMessagePublished("test.subject", 100, time.Millisecond)
		metrics.RecordMessageReceived("test.subject", 100, time.Millisecond)
		metrics.RecordPublishError()
		metrics.RecordReceiveError()
		metrics.RecordRequestSent(100)
		metrics.RecordReplyReceived(100, time.Millisecond)
		metrics.RecordRequestTimeout()
		metrics.RecordSubscriptionCreated()
		metrics.RecordSubscriptionClosed()
		metrics.RecordStreamMessageAdded()
		metrics.RecordStreamMessageDeleted()
		metrics.RecordConsumerMessageAcked()
		metrics.RecordConsumerMessageNacked()
		metrics.RecordConsumerRedelivery()
	})
}

func TestErrorHandling(t *testing.T) {
	// Test error creation and wrapping
	baseErr := assert.AnError
	wrappedErr := natslib.NewNatsError(1001, "test error", baseErr)

	assert.Equal(t, 1001, wrappedErr.Code)
	assert.Equal(t, "test error", wrappedErr.Message)
	assert.Equal(t, baseErr, wrappedErr.Unwrap())
	assert.Contains(t, wrappedErr.Error(), "test error")

	// Test predefined errors
	assert.Equal(t, 1001, natslib.ErrNoVUState.Code)
	assert.Equal(t, 1002, natslib.ErrConnectionClosed.Code)
	assert.Equal(t, 1003, natslib.ErrInvalidConfig.Code)
	assert.Equal(t, 1004, natslib.ErrStreamNotFound.Code)
	assert.Equal(t, 1005, natslib.ErrConsumerNotFound.Code)
	assert.Equal(t, 1006, natslib.ErrTimeout.Code)
	assert.Equal(t, 1007, natslib.ErrNoMessage.Code)
}

func TestUtilityFunctions(t *testing.T) {
	// Test ParseDuration
	assert.Equal(t, 30*time.Second, natslib.ParseDuration(30))
	assert.Equal(t, 0*time.Second, natslib.ParseDuration(0))
	assert.Equal(t, 3600*time.Second, natslib.ParseDuration(3600))

	// Test ParseTimestamp
	ts := int64(1634567890)
	expected := time.Unix(ts, 0)
	assert.Equal(t, expected, natslib.ParseTimestamp(ts))

	// Test zero timestamp
	zero := int64(0)
	expectedZero := time.Unix(zero, 0)
	assert.Equal(t, expectedZero, natslib.ParseTimestamp(zero))
}

func TestConfigurationEdgeCases(t *testing.T) {
	t.Run("empty connection options", func(t *testing.T) {
		opts := natslib.ConnectionOptions{}
		err := natslib.ValidateConnectionOptions(opts)
		assert.NoError(t, err) // Should be valid with defaults
	})

	t.Run("minimal valid stream config", func(t *testing.T) {
		config := natslib.StreamConfig{
			Name:     "TEST",
			Subjects: []string{"test"},
			Replicas: 1,
		}
		err := natslib.ValidateStreamConfig(config)
		assert.NoError(t, err)
	})

	t.Run("minimal valid consumer config", func(t *testing.T) {
		config := natslib.ConsumerConfig{
			Stream: "TEST",
		}
		err := natslib.ValidateConsumerConfig(config)
		assert.NoError(t, err)
	})

	t.Run("TLS options validation", func(t *testing.T) {
		tests := []struct {
			name    string
			tlsOpts *natslib.TLSOptions
			wantErr bool
		}{
			{
				name:    "nil TLS options",
				tlsOpts: nil,
				wantErr: false,
			},
			{
				name: "both cert and key",
				tlsOpts: &natslib.TLSOptions{
					CertFile: "cert.pem",
					KeyFile:  "key.pem",
				},
				wantErr: false,
			},
			{
				name: "cert without key",
				tlsOpts: &natslib.TLSOptions{
					CertFile: "cert.pem",
				},
				wantErr: true,
			},
			{
				name: "key without cert",
				tlsOpts: &natslib.TLSOptions{
					KeyFile: "key.pem",
				},
				wantErr: true,
			},
			{
				name: "CA file only",
				tlsOpts: &natslib.TLSOptions{
					CAFile: "ca.pem",
				},
				wantErr: false,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				opts := natslib.ConnectionOptions{
					TLS: tt.tlsOpts,
				}
				err := natslib.ValidateConnectionOptions(opts)
				if tt.wantErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
			})
		}
	})
}
