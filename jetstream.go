package nats

import (
	"time"

	"github.com/nats-io/nats.go"
)

type StreamConfig struct {
	Name      string   `js:"name"`
	Subjects  []string `js:"subjects"`
	Retention string   `js:"retention"`
	MaxBytes  int64    `js:"maxBytes"`
	MaxMsgs   int64    `js:"maxMsgs"`
	MaxAge    int      `js:"maxAge"`
	Replicas  int      `js:"replicas"`
	Discard   string   `js:"discard"`
	Storage   string   `js:"storage"`
}

func (j *JetStream) AddStream(config StreamConfig) error {
	if j.js == nil {
		return ErrConnectionClosed
	}

	if config.Name == "" {
		return NewNatsError(1015, "stream name cannot be empty", nil)
	}

	// Convert retention policy
	var retention nats.RetentionPolicy
	switch config.Retention {
	case "limits":
		retention = nats.LimitsPolicy
	case "interest":
		retention = nats.InterestPolicy
	case "workqueue":
		retention = nats.WorkQueuePolicy
	default:
		retention = nats.LimitsPolicy
	}

	// Convert discard policy
	var discard nats.DiscardPolicy
	switch config.Discard {
	case "old":
		discard = nats.DiscardOld
	case "new":
		discard = nats.DiscardNew
	default:
		discard = nats.DiscardOld
	}

	// Convert storage type
	var storage nats.StorageType
	switch config.Storage {
	case "file":
		storage = nats.FileStorage
	case "memory":
		storage = nats.MemoryStorage
	default:
		storage = nats.FileStorage
	}

	streamConfig := &nats.StreamConfig{
		Name:      config.Name,
		Subjects:  config.Subjects,
		Retention: retention,
		MaxBytes:  config.MaxBytes,
		MaxMsgs:   config.MaxMsgs,
		MaxAge:    time.Duration(config.MaxAge) * time.Second,
		Replicas:  config.Replicas,
		Discard:   discard,
		Storage:   storage,
	}

	_, err := j.js.AddStream(streamConfig)
	if err != nil {
		return NewNatsError(1016, "failed to add stream", err)
	}

	return nil
}

func (j *JetStream) UpdateStream(config StreamConfig) error {
	if j.js == nil {
		return ErrConnectionClosed
	}

	if config.Name == "" {
		return NewNatsError(1015, "stream name cannot be empty", nil)
	}

	// Get existing stream info first
	info, err := j.js.StreamInfo(config.Name)
	if err != nil {
		return NewNatsError(1017, "stream not found", err)
	}

	// Update existing config
	streamConfig := info.Config

	if len(config.Subjects) > 0 {
		streamConfig.Subjects = config.Subjects
	}
	if config.MaxBytes > 0 {
		streamConfig.MaxBytes = config.MaxBytes
	}
	if config.MaxMsgs > 0 {
		streamConfig.MaxMsgs = config.MaxMsgs
	}
	if config.MaxAge > 0 {
		streamConfig.MaxAge = time.Duration(config.MaxAge) * time.Second
	}
	if config.Replicas > 0 {
		streamConfig.Replicas = config.Replicas
	}

	_, err = j.js.UpdateStream(&streamConfig)
	if err != nil {
		return NewNatsError(1018, "failed to update stream", err)
	}

	return nil
}

func (j *JetStream) DeleteStream(streamName string) error {
	if j.js == nil {
		return ErrConnectionClosed
	}

	if streamName == "" {
		return NewNatsError(1015, "stream name cannot be empty", nil)
	}

	err := j.js.DeleteStream(streamName)
	if err != nil {
		return NewNatsError(1019, "failed to delete stream", err)
	}

	return nil
}

func (j *JetStream) StreamInfo(streamName string) (*nats.StreamInfo, error) {
	if j.js == nil {
		return nil, ErrConnectionClosed
	}

	if streamName == "" {
		return nil, NewNatsError(1015, "stream name cannot be empty", nil)
	}

	info, err := j.js.StreamInfo(streamName)
	if err != nil {
		return nil, NewNatsError(1020, "failed to get stream info", err)
	}

	return info, nil
}

func (j *JetStream) ListStreams() ([]string, error) {
	if j.js == nil {
		return nil, ErrConnectionClosed
	}

	streamNames := j.js.StreamNames()
	var streams []string
	for name := range streamNames {
		streams = append(streams, name)
	}

	return streams, nil
}

func (j *JetStream) Publish(subject string, data []byte) error {
	if j.js == nil {
		return ErrConnectionClosed
	}

	if subject == "" {
		return NewNatsError(1008, "subject cannot be empty", nil)
	}

	_, err := j.js.Publish(subject, data)
	if err != nil {
		return NewNatsError(1021, "failed to publish to jetstream", err)
	}

	return nil
}

func (j *JetStream) PublishAsync(subject string, data []byte) error {
	if j.js == nil {
		return ErrConnectionClosed
	}

	if subject == "" {
		return NewNatsError(1008, "subject cannot be empty", nil)
	}

	_, err := j.js.PublishAsync(subject, data)
	if err != nil {
		return NewNatsError(1022, "failed to publish async to jetstream", err)
	}

	return nil
}

// GetStreamInfo retrieves detailed information about a stream
func (j *JetStream) GetStreamInfo(streamName string) (*nats.StreamInfo, error) {
	if j.js == nil {
		return nil, ErrConnectionClosed
	}

	if streamName == "" {
		return nil, NewNatsError(1015, "stream name cannot be empty", nil)
	}

	info, err := j.js.StreamInfo(streamName)
	if err != nil {
		return nil, NewNatsError(1020, "failed to get stream info", err)
	}

	return info, nil
}

// GetConsumerInfo retrieves detailed information about a consumer
func (j *JetStream) GetConsumerInfo(streamName, consumerName string) (*nats.ConsumerInfo, error) {
	if j.js == nil {
		return nil, ErrConnectionClosed
	}

	if streamName == "" {
		return nil, NewNatsError(1015, "stream name cannot be empty", nil)
	}

	if consumerName == "" {
		return nil, NewNatsError(1024, "consumer name cannot be empty", nil)
	}

	info, err := j.js.ConsumerInfo(streamName, consumerName)
	if err != nil {
		return nil, NewNatsError(1029, "failed to get consumer info", err)
	}

	return info, nil
}

// GetAccountInfo retrieves JetStream account information
func (j *JetStream) GetAccountInfo() (*nats.AccountInfo, error) {
	if j.js == nil {
		return nil, ErrConnectionClosed
	}

	info, err := j.js.AccountInfo()
	if err != nil {
		return nil, NewNatsError(1034, "failed to get account info", err)
	}

	return info, nil
}

// PurgeStream removes all messages from a stream
func (j *JetStream) PurgeStream(streamName string) error {
	if j.js == nil {
		return ErrConnectionClosed
	}

	if streamName == "" {
		return NewNatsError(1015, "stream name cannot be empty", nil)
	}

	err := j.js.PurgeStream(streamName)
	if err != nil {
		return NewNatsError(1035, "failed to purge stream", err)
	}

	return nil
}

// DeleteMessage removes a specific message from a stream
func (j *JetStream) DeleteMessage(streamName string, seq uint64) error {
	if j.js == nil {
		return ErrConnectionClosed
	}

	if streamName == "" {
		return NewNatsError(1015, "stream name cannot be empty", nil)
	}

	// Use direct API call to delete message
	_, err := j.js.StreamInfo(streamName)
	if err != nil {
		return NewNatsError(1020, "stream not found", err)
	}

	// Note: NATS Go client doesn't expose direct message deletion
	// This would need to be implemented via direct API calls
	return NewNatsError(1036, "message deletion not implemented in current NATS client", nil)
}

// GetStreamNames returns all stream names
func (j *JetStream) GetStreamNames() ([]string, error) {
	if j.js == nil {
		return nil, ErrConnectionClosed
	}

	streamNames := j.js.StreamNames()
	var streams []string
	for name := range streamNames {
		streams = append(streams, name)
	}

	return streams, nil
}

// GetConsumerNames returns all consumer names for a stream
func (j *JetStream) GetConsumerNames(streamName string) ([]string, error) {
	if j.js == nil {
		return nil, ErrConnectionClosed
	}

	if streamName == "" {
		return nil, NewNatsError(1015, "stream name cannot be empty", nil)
	}

	consumerNames := j.js.ConsumerNames(streamName)
	var consumers []string
	for name := range consumerNames {
		consumers = append(consumers, name)
	}

	return consumers, nil
}
