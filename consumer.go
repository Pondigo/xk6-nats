package nats

import (
	"time"

	"github.com/nats-io/nats.go"
)

type ConsumerConfig struct {
	Stream        string `js:"stream"`
	Name          string `js:"name"`
	Durable       string `js:"durable"`
	DeliverPolicy string `js:"deliverPolicy"`
	OptStartSeq   uint64 `js:"optStartSeq"`
	OptStartTime  int64  `js:"optStartTime"`
	AckPolicy     string `js:"ackPolicy"`
	AckWait       int    `js:"ackWait"`
	MaxDeliver    int    `js:"maxDeliver"`
	BackOff       []int  `js:"backOff"`
	FilterSubject string `js:"filterSubject"`
	ReplayPolicy  string `js:"replayPolicy"`
	SampleFreq    string `js:"sampleFreq"`
}

func (j *JetStream) AddConsumer(streamName string, config ConsumerConfig) error {
	if j.js == nil {
		return ErrConnectionClosed
	}

	if streamName == "" {
		return NewNatsError(1015, "stream name cannot be empty", nil)
	}

	if config.Durable == "" {
		return NewNatsError(1024, "consumer durable name cannot be empty", nil)
	}

	// Convert deliver policy
	var deliverPolicy nats.DeliverPolicy
	switch config.DeliverPolicy {
	case "all":
		deliverPolicy = nats.DeliverAllPolicy
	case "last":
		deliverPolicy = nats.DeliverLastPolicy
	case "new":
		deliverPolicy = nats.DeliverNewPolicy
	case "by_start_sequence":
		deliverPolicy = nats.DeliverByStartSequencePolicy
	case "by_start_time":
		deliverPolicy = nats.DeliverByStartTimePolicy
	case "last_per_subject":
		deliverPolicy = nats.DeliverLastPerSubjectPolicy
	default:
		deliverPolicy = nats.DeliverAllPolicy
	}

	// Convert ack policy
	var ackPolicy nats.AckPolicy
	switch config.AckPolicy {
	case "none":
		ackPolicy = nats.AckNonePolicy
	case "all":
		ackPolicy = nats.AckAllPolicy
	case "explicit":
		ackPolicy = nats.AckExplicitPolicy
	default:
		ackPolicy = nats.AckExplicitPolicy
	}

	// Convert replay policy
	var replayPolicy nats.ReplayPolicy
	switch config.ReplayPolicy {
	case "instant":
		replayPolicy = nats.ReplayInstantPolicy
	case "original":
		replayPolicy = nats.ReplayOriginalPolicy
	default:
		replayPolicy = nats.ReplayInstantPolicy
	}

	consumerConfig := &nats.ConsumerConfig{
		Durable:       config.Durable,
		DeliverPolicy: deliverPolicy,
		OptStartSeq:   config.OptStartSeq,
		AckPolicy:     ackPolicy,
		AckWait:       time.Duration(config.AckWait) * time.Second,
		MaxDeliver:    config.MaxDeliver,
		FilterSubject: config.FilterSubject,
		ReplayPolicy:  replayPolicy,
	}

	if config.OptStartTime > 0 {
		startTime := time.Unix(config.OptStartTime, 0)
		consumerConfig.OptStartTime = &startTime
	}

	if len(config.BackOff) > 0 {
		for _, backoff := range config.BackOff {
			consumerConfig.BackOff = append(consumerConfig.BackOff, time.Duration(backoff)*time.Second)
		}
	}

	_, err := j.js.AddConsumer(streamName, consumerConfig)
	if err != nil {
		return NewNatsError(1025, "failed to add consumer", err)
	}

	return nil
}

func (j *JetStream) UpdateConsumer(streamName string, config ConsumerConfig) error {
	if j.js == nil {
		return ErrConnectionClosed
	}

	if streamName == "" {
		return NewNatsError(1015, "stream name cannot be empty", nil)
	}

	if config.Durable == "" {
		return NewNatsError(1024, "consumer durable name cannot be empty", nil)
	}

	// Get existing consumer info first
	info, err := j.js.ConsumerInfo(streamName, config.Durable)
	if err != nil {
		return NewNatsError(1026, "consumer not found", err)
	}

	// Update existing config
	consumerConfig := info.Config

	if config.FilterSubject != "" {
		consumerConfig.FilterSubject = config.FilterSubject
	}
	if config.AckWait > 0 {
		consumerConfig.AckWait = time.Duration(config.AckWait) * time.Second
	}
	if config.MaxDeliver > 0 {
		consumerConfig.MaxDeliver = config.MaxDeliver
	}

	_, err = j.js.UpdateConsumer(streamName, &consumerConfig)
	if err != nil {
		return NewNatsError(1027, "failed to update consumer", err)
	}

	return nil
}

func (j *JetStream) DeleteConsumer(streamName, consumerName string) error {
	if j.js == nil {
		return ErrConnectionClosed
	}

	if streamName == "" {
		return NewNatsError(1015, "stream name cannot be empty", nil)
	}

	if consumerName == "" {
		return NewNatsError(1024, "consumer name cannot be empty", nil)
	}

	err := j.js.DeleteConsumer(streamName, consumerName)
	if err != nil {
		return NewNatsError(1028, "failed to delete consumer", err)
	}

	return nil
}

func (j *JetStream) ConsumerInfo(streamName, consumerName string) (*nats.ConsumerInfo, error) {
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

func (j *JetStream) PullSubscribe(streamName, subject, durable string) (*nats.Subscription, error) {
	if j.js == nil {
		return nil, ErrConnectionClosed
	}

	if subject == "" {
		return nil, NewNatsError(1008, "subject cannot be empty", nil)
	}

	if durable == "" {
		return nil, NewNatsError(1024, "durable name cannot be empty", nil)
	}

	sub, err := j.js.PullSubscribe(subject, durable)
	if err != nil {
		return nil, NewNatsError(1030, "failed to create pull subscription", err)
	}

	return sub, nil
}

func (j *JetStream) PullMessages(sub *nats.Subscription, batchSize int, timeout time.Duration) ([]*nats.Msg, error) {
	if sub == nil {
		return nil, NewNatsError(1031, "subscription cannot be nil", nil)
	}

	if batchSize <= 0 {
		batchSize = 1
	}

	if timeout <= 0 {
		timeout = 30 * time.Second
	}

	msgs, err := sub.Fetch(batchSize, nats.MaxWait(timeout))
	if err != nil && err != nats.ErrTimeout {
		return nil, NewNatsError(1032, "failed to fetch messages", err)
	}

	return msgs, nil
}

func (j *JetStream) PushSubscribe(streamName, subject, durable string, handler func(*nats.Msg)) (*nats.Subscription, error) {
	if j.js == nil {
		return nil, ErrConnectionClosed
	}

	if subject == "" {
		return nil, NewNatsError(1008, "subject cannot be empty", nil)
	}

	natsHandler := func(msg *nats.Msg) {
		if j.vu.State() == nil {
			return
		}

		j.vu.State().Logger.Debugf("Received push message on subject %s", msg.Subject)
		handler(msg)
	}

	var sub *nats.Subscription
	var err error

	if durable != "" {
		sub, err = j.js.Subscribe(subject, natsHandler, nats.Durable(durable))
	} else {
		sub, err = j.js.Subscribe(subject, natsHandler)
	}

	if err != nil {
		return nil, NewNatsError(1033, "failed to create push subscription", err)
	}

	return sub, nil
}

func (j *JetStream) ListConsumers(streamName string) ([]string, error) {
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
