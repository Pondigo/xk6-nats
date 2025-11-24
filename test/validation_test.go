package test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	natslib "github.com/pondigo/xk6-nats"
)

func TestConnectionOptions(t *testing.T) {
	tests := []struct {
		name    string
		opts    natslib.ConnectionOptions
		wantErr bool
	}{
		{
			name: "valid options",
			opts: natslib.ConnectionOptions{
				URLs:           []string{"nats://localhost:4222"},
				MaxReconnects:  10,
				ReconnectWait:  2,
				PingInterval:   60,
				MaxPingsOut:    2,
				AllowReconnect: true,
			},
			wantErr: false,
		},
		{
			name:    "empty URLs should be valid",
			opts:    natslib.ConnectionOptions{},
			wantErr: false,
		},
		{
			name: "negative maxReconnects",
			opts: natslib.ConnectionOptions{
				MaxReconnects: -1,
			},
			wantErr: true,
		},
		{
			name: "negative reconnectWait",
			opts: natslib.ConnectionOptions{
				ReconnectWait: -1,
			},
			wantErr: true,
		},
		{
			name: "negative pingInterval",
			opts: natslib.ConnectionOptions{
				PingInterval: -1,
			},
			wantErr: true,
		},
		{
			name: "negative maxPingsOut",
			opts: natslib.ConnectionOptions{
				MaxPingsOut: -1,
			},
			wantErr: true,
		},
		{
			name: "TLS cert without key",
			opts: natslib.ConnectionOptions{
				TLS: &natslib.TLSOptions{
					CertFile: "cert.pem",
				},
			},
			wantErr: true,
		},
		{
			name: "TLS key without cert",
			opts: natslib.ConnectionOptions{
				TLS: &natslib.TLSOptions{
					KeyFile: "key.pem",
				},
			},
			wantErr: true,
		},
		{
			name: "valid TLS options",
			opts: natslib.ConnectionOptions{
				TLS: &natslib.TLSOptions{
					CertFile: "cert.pem",
					KeyFile:  "key.pem",
					CAFile:   "ca.pem",
					Insecure: false,
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := natslib.ValidateConnectionOptions(tt.opts)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestStreamConfig(t *testing.T) {
	tests := []struct {
		name    string
		config  natslib.StreamConfig
		wantErr bool
	}{
		{
			name: "valid stream config",
			config: natslib.StreamConfig{
				Name:     "TEST_STREAM",
				Subjects: []string{"test.>"},
				Replicas: 1,
			},
			wantErr: false,
		},
		{
			name: "empty name",
			config: natslib.StreamConfig{
				Subjects: []string{"test.>"},
			},
			wantErr: true,
		},
		{
			name: "no subjects",
			config: natslib.StreamConfig{
				Name: "TEST_STREAM",
			},
			wantErr: true,
		},
		{
			name: "invalid replicas - too low",
			config: natslib.StreamConfig{
				Name:     "TEST_STREAM",
				Subjects: []string{"test.>"},
				Replicas: 0,
			},
			wantErr: true,
		},
		{
			name: "invalid replicas - too high",
			config: natslib.StreamConfig{
				Name:     "TEST_STREAM",
				Subjects: []string{"test.>"},
				Replicas: 6,
			},
			wantErr: true,
		},
		{
			name: "negative maxBytes",
			config: natslib.StreamConfig{
				Name:     "TEST_STREAM",
				Subjects: []string{"test.>"},
				MaxBytes: -1,
			},
			wantErr: true,
		},
		{
			name: "negative maxMsgs",
			config: natslib.StreamConfig{
				Name:     "TEST_STREAM",
				Subjects: []string{"test.>"},
				MaxMsgs:  -1,
			},
			wantErr: true,
		},
		{
			name: "negative maxAge",
			config: natslib.StreamConfig{
				Name:     "TEST_STREAM",
				Subjects: []string{"test.>"},
				MaxAge:   -1,
			},
			wantErr: true,
		},
		{
			name: "valid retention policies",
			config: natslib.StreamConfig{
				Name:      "TEST_STREAM",
				Subjects:  []string{"test.>"},
				Retention: "limits",
				Replicas:  1,
			},
			wantErr: false,
		},
		{
			name: "valid storage types - duplicate",
			config: natslib.StreamConfig{
				Name:     "TEST_STREAM",
				Subjects: []string{"test.>"},
				Storage:  "file",
				Replicas: 1,
			},
			wantErr: false,
		},
		{
			name: "valid discard policies - duplicate",
			config: natslib.StreamConfig{
				Name:     "TEST_STREAM",
				Subjects: []string{"test.>"},
				Discard:  "old",
				Replicas: 1,
			},
			wantErr: false,
		},
		{
			name: "valid storage types",
			config: natslib.StreamConfig{
				Name:     "TEST_STREAM",
				Subjects: []string{"test.>"},
				Storage:  "file",
				Replicas: 1,
			},
			wantErr: false,
		},
		{
			name: "valid discard policies",
			config: natslib.StreamConfig{
				Name:     "TEST_STREAM",
				Subjects: []string{"test.>"},
				Discard:  "old",
				Replicas: 1,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := natslib.ValidateStreamConfig(tt.config)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestConsumerConfig(t *testing.T) {
	tests := []struct {
		name    string
		config  natslib.ConsumerConfig
		wantErr bool
	}{
		{
			name: "valid consumer config",
			config: natslib.ConsumerConfig{
				Stream:     "TEST_STREAM",
				Durable:    "TEST_CONSUMER",
				AckWait:    30,
				MaxDeliver: 3,
			},
			wantErr: false,
		},
		{
			name: "empty stream",
			config: natslib.ConsumerConfig{
				Durable: "TEST_CONSUMER",
				AckWait: 30,
			},
			wantErr: true,
		},
		{
			name: "negative ackWait",
			config: natslib.ConsumerConfig{
				Stream:  "TEST_STREAM",
				AckWait: -1,
			},
			wantErr: true,
		},
		{
			name: "negative maxDeliver",
			config: natslib.ConsumerConfig{
				Stream:     "TEST_STREAM",
				MaxDeliver: -1,
			},
			wantErr: true,
		},
		{
			name: "negative backoff",
			config: natslib.ConsumerConfig{
				Stream:  "TEST_STREAM",
				BackOff: []int{-1, 2, 3},
			},
			wantErr: true,
		},
		{
			name: "valid deliver policies",
			config: natslib.ConsumerConfig{
				Stream:        "TEST_STREAM",
				DeliverPolicy: "all",
			},
			wantErr: false,
		},
		{
			name: "valid ack policies",
			config: natslib.ConsumerConfig{
				Stream:    "TEST_STREAM",
				AckPolicy: "explicit",
			},
			wantErr: false,
		},
		{
			name: "valid replay policies",
			config: natslib.ConsumerConfig{
				Stream:       "TEST_STREAM",
				ReplayPolicy: "instant",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := natslib.ValidateConsumerConfig(tt.config)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestParseDuration(t *testing.T) {
	tests := []struct {
		seconds  int
		expected time.Duration
	}{
		{0, 0 * time.Second},
		{1, 1 * time.Second},
		{30, 30 * time.Second},
		{60, 60 * time.Second},
		{3600, 3600 * time.Second},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := natslib.ParseDuration(tt.seconds)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestParseTimestamp(t *testing.T) {
	tests := []struct {
		timestamp int64
		expected  time.Time
	}{
		{0, time.Unix(0, 0)},
		{1634567890, time.Unix(1634567890, 0)},
		{1640995200, time.Unix(1640995200, 0)}, // 2022-01-01 00:00:00 UTC
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := natslib.ParseTimestamp(tt.timestamp)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNatsError(t *testing.T) {
	// Test error creation
	err := natslib.NewNatsError(1001, "test error", nil)
	assert.Equal(t, 1001, err.Code)
	assert.Equal(t, "test error", err.Message)
	assert.Nil(t, err.Unwrap())
	assert.Contains(t, err.Error(), "nats error [1001]: test error")

	// Test error with underlying error
	underlyingErr := natslib.NewNatsError(999, "underlying", nil)
	err = natslib.NewNatsError(1002, "wrapper error", underlyingErr)
	assert.Equal(t, 1002, err.Code)
	assert.Equal(t, "wrapper error", err.Message)
	assert.Equal(t, underlyingErr, err.Unwrap())
	assert.Contains(t, err.Error(), "nats error [1002]: wrapper error: nats error [999]: underlying")
}

func TestPredefinedErrors(t *testing.T) {
	tests := []struct {
		name         string
		err          *natslib.NatsError
		expectedCode int
		expectedMsg  string
	}{
		{"NoVUState", natslib.ErrNoVUState, 1001, "no VU state available"},
		{"ConnectionClosed", natslib.ErrConnectionClosed, 1002, "connection is closed"},
		{"InvalidConfig", natslib.ErrInvalidConfig, 1003, "invalid configuration"},
		{"StreamNotFound", natslib.ErrStreamNotFound, 1004, "stream not found"},
		{"ConsumerNotFound", natslib.ErrConsumerNotFound, 1005, "consumer not found"},
		{"Timeout", natslib.ErrTimeout, 1006, "operation timed out"},
		{"NoMessage", natslib.ErrNoMessage, 1007, "no message available"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expectedCode, tt.err.Code)
			assert.Equal(t, tt.expectedMsg, tt.err.Message)
			assert.Contains(t, tt.err.Error(), tt.expectedMsg)
		})
	}
}

func TestConnectionError(t *testing.T) {
	// Test connection error creation
	err := natslib.NewConnectionError("connection failed", nil)
	assert.Equal(t, 1002, err.Code)
	assert.Equal(t, "connection failed", err.Message)
	assert.Contains(t, err.Error(), "nats error [1002]: connection failed")

	// Test connection error with underlying error
	underlyingErr := assert.AnError
	err = natslib.NewConnectionError("connection failed", underlyingErr)
	assert.Equal(t, 1002, err.Code)
	assert.Equal(t, "connection failed", err.Message)
	assert.Equal(t, underlyingErr, err.Unwrap())
}
