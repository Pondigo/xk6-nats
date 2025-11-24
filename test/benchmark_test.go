package test

import (
	"testing"
	"time"

	natslib "github.com/pondigo/xk6-nats"
)

func BenchmarkConnectionValidation(b *testing.B) {
	opts := natslib.ConnectionOptions{
		URLs:           []string{"nats://localhost:4222"},
		MaxReconnects:  10,
		ReconnectWait:  2,
		PingInterval:   60,
		MaxPingsOut:    2,
		AllowReconnect: true,
		TLS: &natslib.TLSOptions{
			CertFile: "cert.pem",
			KeyFile:  "key.pem",
			CAFile:   "ca.pem",
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = natslib.ValidateConnectionOptions(opts)
	}
}

func BenchmarkStreamConfigValidation(b *testing.B) {
	config := natslib.StreamConfig{
		Name:      "TEST_STREAM",
		Subjects:  []string{"test.>", "foo.>", "bar.>"},
		Retention: "limits",
		Storage:   "file",
		Discard:   "old",
		Replicas:  3,
		MaxBytes:  1024 * 1024 * 1024, // 1GB
		MaxMsgs:   1000000,
		MaxAge:    86400, // 24 hours
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = natslib.ValidateStreamConfig(config)
	}
}

func BenchmarkConsumerConfigValidation(b *testing.B) {
	config := natslib.ConsumerConfig{
		Stream:        "TEST_STREAM",
		Durable:       "TEST_CONSUMER",
		DeliverPolicy: "all",
		AckPolicy:     "explicit",
		AckWait:       30,
		MaxDeliver:    3,
		BackOff:       []int{1, 2, 4, 8, 16},
		FilterSubject: "test.>",
		ReplayPolicy:  "instant",
		SampleFreq:    "100%",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = natslib.ValidateConsumerConfig(config)
	}
}

func BenchmarkErrorCreation(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = natslib.NewNatsError(1001, "test error", nil)
	}
}

func BenchmarkParseDuration(b *testing.B) {
	seconds := []int{0, 1, 30, 60, 300, 3600, 86400}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = natslib.ParseDuration(seconds[i%len(seconds)])
	}
}

func BenchmarkParseTimestamp(b *testing.B) {
	timestamps := []int64{0, 1634567890, 1640995200, 1672531200, 1704067200}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = natslib.ParseTimestamp(timestamps[i%len(timestamps)])
	}
}

func BenchmarkMetricsCreation(b *testing.B) {
	vu := &MockVU{
		ctx: nil,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = natslib.NewNatsMetrics(vu)
	}
}

func BenchmarkMetricsRecording(b *testing.B) {
	vu := &MockVU{
		ctx: nil,
	}

	metrics, _ := natslib.NewNatsMetrics(vu)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		switch i % 10 {
		case 0:
			metrics.RecordConnectionEstablished()
		case 1:
			metrics.RecordMessagePublished("test.subject", 1024, time.Millisecond)
		case 2:
			metrics.RecordMessageReceived("test.subject", 512, 2*time.Millisecond)
		case 3:
			metrics.RecordRequestSent(256)
		case 4:
			metrics.RecordReplyReceived(128, 5*time.Millisecond)
		case 5:
			metrics.RecordSubscriptionCreated()
		case 6:
			metrics.RecordStreamMessageAdded()
		case 7:
			metrics.RecordConsumerMessageAcked()
		case 8:
			metrics.RecordPublishError()
		case 9:
			metrics.RecordConnectionClosed()
		}
	}
}

// BenchmarkComplexValidation tests validation with complex configurations
func BenchmarkComplexValidation(b *testing.B) {
	// Complex connection options
	connOpts := natslib.ConnectionOptions{
		URLs:           []string{"nats://server1:4222", "nats://server2:4222", "nats://server3:4222"},
		MaxReconnects:  100,
		ReconnectWait:  5,
		PingInterval:   30,
		MaxPingsOut:    5,
		AllowReconnect: true,
		User:           "testuser",
		Password:       "testpass",
		Token:          "testtoken",
		TLS: &natslib.TLSOptions{
			CertFile: "client.crt",
			KeyFile:  "client.key",
			CAFile:   "ca.crt",
			Insecure: false,
		},
	}

	// Complex stream config
	streamConfig := natslib.StreamConfig{
		Name:      "COMPLEX_STREAM",
		Subjects:  []string{"events.>", "commands.>", "queries.>"},
		Retention: "limits",
		MaxBytes:  10 * 1024 * 1024 * 1024, // 10GB
		MaxMsgs:   50000000,
		MaxAge:    7 * 86400, // 7 days
		Replicas:  5,
		Discard:   "old",
		Storage:   "file",
	}

	// Complex consumer config
	consumerConfig := natslib.ConsumerConfig{
		Stream:        "COMPLEX_STREAM",
		Durable:       "COMPLEX_CONSUMER",
		DeliverPolicy: "by_start_sequence",
		OptStartSeq:   1000000,
		AckPolicy:     "explicit",
		AckWait:       60,
		MaxDeliver:    10,
		BackOff:       []int{1, 2, 4, 8, 16, 32, 64},
		FilterSubject: "events.important.>",
		ReplayPolicy:  "original",
		SampleFreq:    "0.1%",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		switch i % 3 {
		case 0:
			_ = natslib.ValidateConnectionOptions(connOpts)
		case 1:
			_ = natslib.ValidateStreamConfig(streamConfig)
		case 2:
			_ = natslib.ValidateConsumerConfig(consumerConfig)
		}
	}
}
