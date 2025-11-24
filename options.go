package nats

import (
	"fmt"
	"time"
)

func ValidateConnectionOptions(opts ConnectionOptions) error {
	if len(opts.URLs) == 0 {
		opts.URLs = []string{"nats://localhost:4222"}
	}

	if opts.MaxReconnects < 0 {
		return fmt.Errorf("maxReconnects must be non-negative")
	}

	if opts.ReconnectWait < 0 {
		return fmt.Errorf("reconnectWait must be non-negative")
	}

	if opts.PingInterval < 0 {
		return fmt.Errorf("pingInterval must be non-negative")
	}

	if opts.MaxPingsOut < 0 {
		return fmt.Errorf("maxPingsOut must be non-negative")
	}

	if opts.TLS != nil {
		if opts.TLS.CertFile != "" && opts.TLS.KeyFile == "" {
			return fmt.Errorf("certFile specified but keyFile is missing")
		}
		if opts.TLS.KeyFile != "" && opts.TLS.CertFile == "" {
			return fmt.Errorf("keyFile specified but certFile is missing")
		}
	}

	return nil
}

func ValidateStreamConfig(config StreamConfig) error {
	if config.Name == "" {
		return fmt.Errorf("stream name is required")
	}

	if len(config.Subjects) == 0 {
		return fmt.Errorf("at least one subject is required")
	}

	if config.MaxBytes < 0 {
		return fmt.Errorf("maxBytes must be non-negative")
	}

	if config.MaxMsgs < 0 {
		return fmt.Errorf("maxMsgs must be non-negative")
	}

	if config.MaxAge < 0 {
		return fmt.Errorf("maxAge must be non-negative")
	}

	if config.Replicas < 1 || config.Replicas > 5 {
		return fmt.Errorf("replicas must be between 1 and 5")
	}

	return nil
}

func ValidateConsumerConfig(config ConsumerConfig) error {
	if config.Stream == "" {
		return fmt.Errorf("stream name is required")
	}

	if config.AckWait < 0 {
		return fmt.Errorf("ackWait must be non-negative")
	}

	if config.MaxDeliver < 0 {
		return fmt.Errorf("maxDeliver must be non-negative")
	}

	for i, backoff := range config.BackOff {
		if backoff < 0 {
			return fmt.Errorf("backOff[%d] must be non-negative", i)
		}
	}

	return nil
}

func ParseDuration(seconds int) time.Duration {
	return time.Duration(seconds) * time.Second
}

func ParseTimestamp(timestamp int64) time.Time {
	return time.Unix(timestamp, 0)
}
