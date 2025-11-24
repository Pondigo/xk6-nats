package test

import (
	"context"
	"testing"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/require"

	natslib "github.com/pondigo/xk6-nats"
)

// MockVU implements the VU interface for testing
type MockVU struct {
	state any
	ctx   context.Context
}

func (m *MockVU) State() any {
	return m.state
}

func (m *MockVU) Context() context.Context {
	return m.ctx
}

// TestHelper provides utilities for testing
type TestHelper struct {
	t *testing.T
}

// NewTestHelper creates a new test helper
func NewTestHelper(t *testing.T) *TestHelper {
	return &TestHelper{t: t}
}

// GetValidConnectionOptions returns valid connection options for testing
func (h *TestHelper) GetValidConnectionOptions() natslib.ConnectionOptions {
	return natslib.ConnectionOptions{
		URLs:           []string{nats.DefaultURL},
		MaxReconnects:  5,
		ReconnectWait:  2,
		PingInterval:   60,
		MaxPingsOut:    2,
		AllowReconnect: true,
	}
}

// GetValidStreamConfig returns valid stream config for testing
func (h *TestHelper) GetValidStreamConfig() natslib.StreamConfig {
	return natslib.StreamConfig{
		Name:     "TEST_STREAM",
		Subjects: []string{"test.>"},
		Replicas: 1,
	}
}

// GetValidConsumerConfig returns valid consumer config for testing
func (h *TestHelper) GetValidConsumerConfig() natslib.ConsumerConfig {
	return natslib.ConsumerConfig{
		Stream:  "TEST_STREAM",
		Durable: "TEST_CONSUMER",
		AckWait: 30,
	}
}

// GetValidTLSOptions returns valid TLS options for testing
func (h *TestHelper) GetValidTLSOptions() *natslib.TLSOptions {
	return &natslib.TLSOptions{
		CertFile: "cert.pem",
		KeyFile:  "key.pem",
		CAFile:   "ca.pem",
		Insecure: false,
	}
}

// AssertNoError asserts that error is nil
func (h *TestHelper) AssertNoError(err error) {
	require.NoError(h.t, err)
}

// AssertError asserts that error is not nil
func (h *TestHelper) AssertError(err error) {
	require.Error(h.t, err)
}

// AssertEqual asserts that two values are equal
func (h *TestHelper) AssertEqual(expected, actual any) {
	require.Equal(h.t, expected, actual)
}

// CreateMockVU creates a mock VU for testing
func (h *TestHelper) CreateMockVU() *MockVU {
	return &MockVU{
		state: make(map[string]any),
		ctx:   context.Background(),
	}
}

// WaitForCondition waits for a condition to be true or timeout
func (h *TestHelper) WaitForCondition(condition func() bool, timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if condition() {
			return true
		}
		time.Sleep(10 * time.Millisecond)
	}
	return false
}

// TestDataGenerator provides test data
type TestDataGenerator struct{}

// NewTestDataGenerator creates a new test data generator
func NewTestDataGenerator() *TestDataGenerator {
	return &TestDataGenerator{}
}

// RandomString generates a random string of specified length
func (g *TestDataGenerator) RandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[i%len(charset)]
	}
	return string(result)
}

// RandomSubject generates a random NATS subject
func (g *TestDataGenerator) RandomSubject() string {
	return "test." + g.RandomString(8) + ".>"
}

// RandomStreamName generates a random stream name
func (g *TestDataGenerator) RandomStreamName() string {
	return "TEST_" + g.RandomString(12)
}

// RandomConsumerName generates a random consumer name
func (g *TestDataGenerator) RandomConsumerName() string {
	return "CONSUMER_" + g.RandomString(12)
}

// RandomMessage generates random message data
func (g *TestDataGenerator) RandomMessage(size int) []byte {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, size)
	for i := range result {
		result[i] = charset[i%len(charset)]
	}
	return result
}

// TestScenario represents a test scenario
type TestScenario struct {
	Name        string
	Description string
	Setup       func() error
	Execute     func() error
	Teardown    func() error
	ExpectedErr string
}

// RunTestScenario runs a test scenario
func (h *TestHelper) RunTestScenario(scenario TestScenario) {
	h.t.Run(scenario.Name, func(t *testing.T) {
		// Setup
		if scenario.Setup != nil {
			err := scenario.Setup()
			require.NoError(t, err, "Setup failed")
		}

		// Execute
		err := scenario.Execute()
		if scenario.ExpectedErr != "" {
			require.Error(t, err, "Expected error but got none")
			require.Contains(t, err.Error(), scenario.ExpectedErr, "Error message mismatch")
		} else {
			require.NoError(t, err, "Unexpected error")
		}

		// Teardown
		if scenario.Teardown != nil {
			err := scenario.Teardown()
			require.NoError(t, err, "Teardown failed")
		}
	})
}

// PerformanceTestResult represents performance test results
type PerformanceTestResult struct {
	Operation    string
	Duration     time.Duration
	SuccessCount int
	ErrorCount   int
	Throughput   float64
	AvgLatency   time.Duration
	MinLatency   time.Duration
	MaxLatency   time.Duration
}

// RunPerformanceTest runs a performance test
func (h *TestHelper) RunPerformanceTest(
	operation string,
	iterations int,
	operationFunc func() error,
) PerformanceTestResult {
	var totalLatency time.Duration
	var minLatency = time.Hour
	var maxLatency time.Duration
	successCount := 0
	errorCount := 0

	start := time.Now()

	for range iterations {
		opStart := time.Now()
		err := operationFunc()
		opDuration := time.Since(opStart)

		if err != nil {
			errorCount++
		} else {
			successCount++
			totalLatency += opDuration
			if opDuration < minLatency {
				minLatency = opDuration
			}
			if opDuration > maxLatency {
				maxLatency = opDuration
			}
		}
	}

	totalDuration := time.Since(start)
	avgLatency := time.Duration(0)
	if successCount > 0 {
		avgLatency = totalLatency / time.Duration(successCount)
	}

	throughput := float64(successCount) / totalDuration.Seconds()

	return PerformanceTestResult{
		Operation:    operation,
		Duration:     totalDuration,
		SuccessCount: successCount,
		ErrorCount:   errorCount,
		Throughput:   throughput,
		AvgLatency:   avgLatency,
		MinLatency:   minLatency,
		MaxLatency:   maxLatency,
	}
}

// AssertPerformance asserts performance requirements
func (h *TestHelper) AssertPerformance(result PerformanceTestResult, minThroughput float64, maxLatency time.Duration) {
	h.t.Helper()

	require.GreaterOrEqualf(h.t, result.Throughput, minThroughput,
		"Throughput %f is below minimum %f", result.Throughput, minThroughput)

	require.LessOrEqualf(h.t, result.AvgLatency, maxLatency,
		"Average latency %v exceeds maximum %v", result.AvgLatency, maxLatency)
}
