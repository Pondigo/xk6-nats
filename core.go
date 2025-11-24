package nats

import (
	"time"

	"github.com/nats-io/nats.go"
)

func (c *Connection) Publish(subject string, data []byte) error {
	if c.nc == nil {
		return ErrConnectionClosed
	}

	if subject == "" {
		return NewNatsError(1008, "subject cannot be empty", nil)
	}

	if err := c.nc.Publish(subject, data); err != nil {
		return NewNatsError(1009, "publish failed", err)
	}

	return nil
}

func (c *Connection) Subscribe(subject string, queue string, handler func(*nats.Msg)) (*nats.Subscription, error) {
	if c.nc == nil {
		return nil, ErrConnectionClosed
	}

	if subject == "" {
		return nil, NewNatsError(1008, "subject cannot be empty", nil)
	}

	var sub *nats.Subscription
	var err error

	natsHandler := func(msg *nats.Msg) {
		// Create a safe context for the handler
		if c.vu.State() == nil {
			return
		}

		c.vu.State().Logger.Debugf("Received message on subject %s", msg.Subject)
		handler(msg)
	}

	if queue != "" {
		sub, err = c.nc.QueueSubscribe(subject, queue, natsHandler)
	} else {
		sub, err = c.nc.Subscribe(subject, natsHandler)
	}

	if err != nil {
		return nil, NewNatsError(1010, "subscription failed", err)
	}

	return sub, nil
}

func (c *Connection) Request(subject string, data []byte, timeout time.Duration) (*nats.Msg, error) {
	if c.nc == nil {
		return nil, ErrConnectionClosed
	}

	if subject == "" {
		return nil, NewNatsError(1008, "subject cannot be empty", nil)
	}

	if timeout <= 0 {
		timeout = 30 * time.Second // Default timeout
	}

	msg, err := c.nc.Request(subject, data, timeout)
	if err != nil {
		return nil, NewNatsError(1011, "request failed", err)
	}

	return msg, nil
}

func (c *Connection) Drain() error {
	if c.nc == nil {
		return ErrConnectionClosed
	}

	if err := c.nc.Drain(); err != nil {
		return NewNatsError(1012, "drain failed", err)
	}

	return nil
}

func (c *Connection) Flush() error {
	if c.nc == nil {
		return ErrConnectionClosed
	}

	if err := c.nc.Flush(); err != nil {
		return NewNatsError(1013, "flush failed", err)
	}

	return nil
}

func (c *Connection) FlushTimeout(timeout time.Duration) error {
	if c.nc == nil {
		return ErrConnectionClosed
	}

	if err := c.nc.FlushTimeout(timeout); err != nil {
		return NewNatsError(1013, "flush timeout failed", err)
	}

	return nil
}

func (c *Connection) JetStream() (*JetStream, error) {
	if c.nc == nil {
		return nil, ErrConnectionClosed
	}

	js, err := c.nc.JetStream()
	if err != nil {
		return nil, NewNatsError(1014, "jetstream not available", err)
	}

	return &JetStream{
		vu: c.vu,
		js: js,
	}, nil
}
