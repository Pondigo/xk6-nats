package nats

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestValidateConnectionOptions(t *testing.T) {
	tests := []struct {
		name    string
		opts    ConnectionOptions
		wantErr bool
	}{
		{
			name: "valid options",
			opts: ConnectionOptions{
				URLs:          []string{"nats://localhost:4222"},
				MaxReconnects: 10,
				ReconnectWait: 2,
				PingInterval:  60,
				MaxPingsOut:   2,
			},
			wantErr: false,
		},
		{
			name:    "empty URLs should be valid (uses default)",
			opts:    ConnectionOptions{},
			wantErr: false,
		},
		{
			name: "negative maxReconnects",
			opts: ConnectionOptions{
				MaxReconnects: -1,
			},
			wantErr: true,
		},
		{
			name: "negative reconnectWait",
			opts: ConnectionOptions{
				ReconnectWait: -1,
			},
			wantErr: true,
		},
		{
			name: "TLS cert without key",
			opts: ConnectionOptions{
				TLS: &TLSOptions{
					CertFile: "cert.pem",
				},
			},
			wantErr: true,
		},
		{
			name: "TLS key without cert",
			opts: ConnectionOptions{
				TLS: &TLSOptions{
					KeyFile: "key.pem",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateConnectionOptions(tt.opts)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateStreamConfig(t *testing.T) {
	tests := []struct {
		name    string
		config  StreamConfig
		wantErr bool
	}{
		{
			name: "valid stream config",
			config: StreamConfig{
				Name:     "TEST_STREAM",
				Subjects: []string{"test.>"},
				Replicas: 1,
			},
			wantErr: false,
		},
		{
			name: "empty name",
			config: StreamConfig{
				Subjects: []string{"test.>"},
			},
			wantErr: true,
		},
		{
			name: "no subjects",
			config: StreamConfig{
				Name: "TEST_STREAM",
			},
			wantErr: true,
		},
		{
			name: "invalid replicas",
			config: StreamConfig{
				Name:     "TEST_STREAM",
				Subjects: []string{"test.>"},
				Replicas: 6,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateStreamConfig(tt.config)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateConsumerConfig(t *testing.T) {
	tests := []struct {
		name    string
		config  ConsumerConfig
		wantErr bool
	}{
		{
			name: "valid consumer config",
			config: ConsumerConfig{
				Stream:     "TEST_STREAM",
				AckWait:    30,
				MaxDeliver: 3,
			},
			wantErr: false,
		},
		{
			name: "empty stream",
			config: ConsumerConfig{
				AckWait: 30,
			},
			wantErr: true,
		},
		{
			name: "negative ackWait",
			config: ConsumerConfig{
				Stream:  "TEST_STREAM",
				AckWait: -1,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateConsumerConfig(tt.config)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestParseDuration(t *testing.T) {
	assert.Equal(t, 30*time.Second, ParseDuration(30))
	assert.Equal(t, 0*time.Second, ParseDuration(0))
}

func TestParseTimestamp(t *testing.T) {
	ts := int64(1634567890)
	expected := time.Unix(ts, 0)
	assert.Equal(t, expected, ParseTimestamp(ts))
}
