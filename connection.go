package nats

import (
	"crypto/tls"
	"strings"
	"time"

	"github.com/nats-io/nats.go"
)

type ConnectionOptions struct {
	URLs           []string    `js:"urls"`
	MaxReconnects  int         `js:"maxReconnects"`
	ReconnectWait  int         `js:"reconnectWait"`
	PingInterval   int         `js:"pingInterval"`
	MaxPingsOut    int         `js:"maxPingsOut"`
	AllowReconnect bool        `js:"allowReconnect"`
	TLS            *TLSOptions `js:"tls"`
	User           string      `js:"user"`
	Password       string      `js:"password"`
	Token          string      `js:"token"`
}

type TLSOptions struct {
	CertFile string `js:"certFile"`
	KeyFile  string `js:"keyFile"`
	CAFile   string `js:"caFile"`
	Insecure bool   `js:"insecure"`
}

func (n *NatsInstance) Connect(opts ConnectionOptions) (*Connection, error) {
	// Validate options
	if err := ValidateConnectionOptions(opts); err != nil {
		return nil, err
	}

	// Determine URLs
	urls := opts.URLs
	if len(urls) == 0 {
		urls = []string{nats.DefaultURL}
	}

	// Build NATS options
	natsOpts := []nats.Option{
		nats.ReconnectWait(time.Duration(opts.ReconnectWait) * time.Second),
		nats.MaxReconnects(opts.MaxReconnects),
		nats.PingInterval(time.Duration(opts.PingInterval) * time.Second),
		nats.MaxPingsOutstanding(opts.MaxPingsOut),
		nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
			if err != nil {
				n.vu.State().Logger.Warnf("NATS disconnected: %v", err)
			}
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			n.vu.State().Logger.Infof("NATS reconnected to %v", nc.ConnectedUrl())
		}),
		nats.ErrorHandler(func(nc *nats.Conn, sub *nats.Subscription, err error) {
			n.vu.State().Logger.Errorf("NATS error: %v", err)
		}),
	}

	// Add authentication
	if opts.User != "" && opts.Password != "" {
		natsOpts = append(natsOpts, nats.UserInfo(opts.User, opts.Password))
	} else if opts.Token != "" {
		natsOpts = append(natsOpts, nats.Token(opts.Token))
	}

	// Add TLS configuration
	if opts.TLS != nil {
		tlsConfig := &tls.Config{
			InsecureSkipVerify: opts.TLS.Insecure,
		}

		// Load certificates if provided
		if opts.TLS.CertFile != "" && opts.TLS.KeyFile != "" {
			cert, err := tls.LoadX509KeyPair(opts.TLS.CertFile, opts.TLS.KeyFile)
			if err != nil {
				return nil, NewConnectionError("failed to load TLS certificate: %v", err)
			}
			tlsConfig.Certificates = []tls.Certificate{cert}
		}

		natsOpts = append(natsOpts, nats.Secure(tlsConfig))
	}

	// Connect to NATS
	nc, err := nats.Connect(strings.Join(urls, ","), natsOpts...)
	if err != nil {
		return nil, NewConnectionError("failed to connect to NATS: %v", err)
	}

	// Wait for connection to be established
	if !nc.IsConnected() {
		return nil, NewConnectionError("NATS connection not established", nil)
	}

	return &Connection{
		vu: n.vu,
		nc: nc,
	}, nil
}

func (c *Connection) Close() error {
	if c.nc != nil {
		c.nc.Close()
	}
	return nil
}

func (c *Connection) IsConnected() bool {
	return c.nc != nil && c.nc.IsConnected()
}

func (c *Connection) Stats() nats.Statistics {
	if c.nc != nil {
		return c.nc.Stats()
	}
	return nats.Statistics{}
}
