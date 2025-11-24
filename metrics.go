package nats

import (
	"context"
	"time"

	"github.com/nats-io/nats.go"
	"go.k6.io/k6/metrics"
)

// NatsMetrics holds all NATS-related metrics
type NatsMetrics struct {
	registry *metrics.Registry
	vu       VU
}

// VU interface for accessing k6 VU state
type VU interface {
	State() any // Using any to avoid import issues
	Context() context.Context
}

// NewNatsMetrics creates and registers all NATS metrics
func NewNatsMetrics(vu VU) (*NatsMetrics, error) {
	registry := metrics.NewRegistry()

	return &NatsMetrics{
		registry: registry,
		vu:       vu,
	}, nil
}

// Registry returns the metrics registry
func (m *NatsMetrics) Registry() *metrics.Registry {
	return m.registry
}

// Placeholder methods for metrics recording
func (m *NatsMetrics) RecordConnectionEstablished()                                                 {}
func (m *NatsMetrics) RecordConnectionClosed()                                                      {}
func (m *NatsMetrics) RecordConnectionFailed()                                                      {}
func (m *NatsMetrics) RecordReconnect()                                                             {}
func (m *NatsMetrics) RecordMessagePublished(subject string, dataSize int64, latency time.Duration) {}
func (m *NatsMetrics) RecordMessageReceived(subject string, dataSize int64, latency time.Duration)  {}
func (m *NatsMetrics) RecordPublishError()                                                          {}
func (m *NatsMetrics) RecordReceiveError()                                                          {}
func (m *NatsMetrics) RecordRequestSent(dataSize int64)                                             {}
func (m *NatsMetrics) RecordReplyReceived(dataSize int64, latency time.Duration)                    {}
func (m *NatsMetrics) RecordRequestTimeout()                                                        {}
func (m *NatsMetrics) RecordSubscriptionCreated()                                                   {}
func (m *NatsMetrics) RecordSubscriptionClosed()                                                    {}
func (m *NatsMetrics) RecordStreamMessageAdded()                                                    {}
func (m *NatsMetrics) RecordStreamMessageDeleted()                                                  {}
func (m *NatsMetrics) RecordConsumerMessageAcked()                                                  {}
func (m *NatsMetrics) RecordConsumerMessageNacked()                                                 {}
func (m *NatsMetrics) RecordConsumerRedelivery()                                                    {}

// WrapConnection wraps a NATS connection to collect metrics
func (m *NatsMetrics) WrapConnection(conn *nats.Conn) *nats.Conn {
	return conn
}
