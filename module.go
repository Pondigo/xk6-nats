package nats

import (
	"encoding/json"
	"github.com/dop251/goja"
	"github.com/nats-io/nats.go"
	"go.k6.io/k6/js/modules"
)

type RootModule struct{}

func (*RootModule) NewModuleInstance(vu modules.VU) modules.Instance {
	return &NatsInstance{
		vu: vu,
	}
}

type NatsInstance struct {
	vu modules.VU
}

func (n *NatsInstance) Exports() modules.Exports {
	return modules.Exports{
		Default: n,
		Named: map[string]any{
			"Connection":     n.NewConnection,
			"JetStream":      n.NewJetStream,
			"StreamConfig":   n.NewStreamConfig,
			"ConsumerConfig": n.NewConsumerConfig,
			"TLSOptions":     n.NewTLSOptions,
		},
	}
}

func (n *NatsInstance) ConnectFromJS(opts goja.Value) *Connection {
	var connOpts ConnectionOptions
	if opts != nil && !goja.IsUndefined(opts) && !goja.IsNull(opts) {
		// Convert goja.Value to JSON then to struct
		optsJSON, err := json.Marshal(opts.Export())
		if err != nil {
			n.vu.State().Logger.Errorf("Failed to marshal connection options: %v", err)
			return nil
		}

		if err := json.Unmarshal(optsJSON, &connOpts); err != nil {
			n.vu.State().Logger.Errorf("Failed to parse connection options: %v", err)
			return nil
		}
	}

	// Call the Connect method from connection.go
	conn, err := (*NatsInstance)(n).Connect(connOpts)
	if err != nil {
		n.vu.State().Logger.Errorf("Failed to connect to NATS: %v", err)
		return nil
	}
	return conn
}

func (n *NatsInstance) NewConnection(opts goja.Value) *Connection {
	return n.ConnectFromJS(opts)
}

func (n *NatsInstance) NewJetStream(conn *Connection) *JetStream {
	if conn == nil {
		n.vu.State().Logger.Errorf("Connection cannot be nil")
		return nil
	}

	js, err := conn.JetStream()
	if err != nil {
		n.vu.State().Logger.Errorf("Failed to create JetStream context: %v", err)
		return nil
	}
	return js
}

func (n *NatsInstance) NewStreamConfig(opts goja.Value) *StreamConfig {
	var config StreamConfig
	if opts != nil && !goja.IsUndefined(opts) && !goja.IsNull(opts) {
		optsJSON, err := json.Marshal(opts.Export())
		if err != nil {
			n.vu.State().Logger.Errorf("Failed to marshal stream config: %v", err)
			return nil
		}

		if err := json.Unmarshal(optsJSON, &config); err != nil {
			n.vu.State().Logger.Errorf("Failed to parse stream config: %v", err)
			return nil
		}
	}
	return &config
}

func (n *NatsInstance) NewConsumerConfig(opts goja.Value) *ConsumerConfig {
	var config ConsumerConfig
	if opts != nil && !goja.IsUndefined(opts) && !goja.IsNull(opts) {
		optsJSON, err := json.Marshal(opts.Export())
		if err != nil {
			n.vu.State().Logger.Errorf("Failed to marshal consumer config: %v", err)
			return nil
		}

		if err := json.Unmarshal(optsJSON, &config); err != nil {
			n.vu.State().Logger.Errorf("Failed to parse consumer config: %v", err)
			return nil
		}
	}
	return &config
}

func (n *NatsInstance) NewTLSOptions(opts goja.Value) *TLSOptions {
	var tlsOpts TLSOptions
	if opts != nil && !goja.IsUndefined(opts) && !goja.IsNull(opts) {
		optsJSON, err := json.Marshal(opts.Export())
		if err != nil {
			n.vu.State().Logger.Errorf("Failed to marshal TLS options: %v", err)
			return nil
		}

		if err := json.Unmarshal(optsJSON, &tlsOpts); err != nil {
			n.vu.State().Logger.Errorf("Failed to parse TLS options: %v", err)
			return nil
		}
	}
	return &tlsOpts
}

type Connection struct {
	vu modules.VU
	nc *nats.Conn
}

type JetStream struct {
	vu modules.VU
	js nats.JetStreamContext
}
