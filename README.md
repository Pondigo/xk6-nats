# k6/Pondigo/nats

A k6 extension for NATS testing with JetStream support.

## Features

- ✅ Connection management with TLS and authentication
- ✅ Publish/Subscribe messaging
- ✅ Request/Reply pattern
- ✅ JetStream streams and consumers
- ✅ Pull and push consumer support
- ✅ Stream and consumer monitoring
- ✅ Account information retrieval
- ✅ Stream purging and management
- ✅ Comprehensive metrics collection framework
- ✅ Configuration validation
- ✅ JavaScript API exports
- ✅ Error handling with structured error codes

## Development

### Prerequisites

- Nix with flakes enabled
- Go 1.21+

### Getting Started

```bash
# Enter development environment
nix develop

# Build the extension
go build ./...

# Run tests
go test ./...

# Run linter
golangci-lint run

# Start NATS server for testing
nats-server
```

### Project Structure

```
xk6-nats/
├── module.go          # k6 registration, Exports()
├── connection.go      # Connection management, auth, TLS
├── core.go            # Publish, Subscribe, Request/Reply  
├── jetstream.go       # JetStream context + stream ops
├── consumer.go        # Pull/push consumer handling
├── metrics.go         # k6 metrics registration & emission
├── options.go         # Configuration structs with validation
├── errors.go          # Error types and wrapping
├── flake.nix          # Nix development environment
└── go.mod             # Go module definition
```

### Usage in k6

```javascript
import nats from 'k6/Pondigo/nats';

export default function() {
    // Connect to NATS
    const conn = nats.connect({
        urls: ['nats://localhost:4222'],
        reconnectWait: 2,
        maxReconnects: 10
    });
    
    // Basic messaging
    conn.publish('test.subject', 'hello world');
    
    // Request/Reply
    const reply = conn.request('test.request', 'ping', 5000);
    console.log('Reply:', String(reply.data));
    
    // Subscriptions
    conn.subscribe('test.sub', '', (msg) => {
        console.log('Received:', String(msg.data));
    });
    
    // JetStream operations
    const js = nats.jetStream(conn);
    
    // Create stream
    const streamConfig = nats.streamConfig({
        name: 'TEST_STREAM',
        subjects: ['test.>'],
        retention: 'limits',
        storage: 'file',
        replicas: 1
    });
    js.addStream(streamConfig);
    
    // Create consumer
    const consumerConfig = nats.consumerConfig({
        stream: 'TEST_STREAM',
        durable: 'TEST_CONSUMER',
        deliverPolicy: 'all',
        ackPolicy: 'explicit'
    });
    js.addConsumer('TEST_STREAM', consumerConfig);
    
    // Publish to JetStream
    js.publish('test.js', 'Hello JetStream!');
    
    // Get information
    const streamInfo = js.getStreamInfo('TEST_STREAM');
    const consumerInfo = js.getConsumerInfo('TEST_STREAM', 'TEST_CONSUMER');
    const accountInfo = js.getAccountInfo();
    
    // Pull consumer
    const sub = js.pullSubscribe('TEST_STREAM', 'test.>', 'TEST_CONSUMER');
    const messages = js.pullMessages(sub, 10, 5000);
    
    // Cleanup
    js.deleteConsumer('TEST_STREAM', 'TEST_CONSUMER');
    js.deleteStream('TEST_STREAM');
    conn.close();
}
```

### API Reference

#### Connection Management
- `nats.connect(options)` - Create NATS connection
- `conn.close()` - Close connection
- `conn.is_connected()` - Check connection status
- `conn.stats()` - Get connection statistics

#### Messaging
- `conn.publish(subject, data)` - Publish message
- `conn.subscribe(subject, queue, handler)` - Create subscription
- `conn.request(subject, data, timeout)` - Send request and wait for reply

#### JetStream
- `nats.jetStream(connection)` - Create JetStream context
- `js.addStream(config)` - Create stream
- `js.updateStream(config)` - Update stream
- `js.deleteStream(name)` - Delete stream
- `js.getStreamInfo(name)` - Get stream information
- `js.getStreamNames()` - List all streams
- `js.purgeStream(name)` - Remove all messages from stream

#### Consumers
- `js.addConsumer(stream, config)` - Create consumer
- `js.updateConsumer(stream, config)` - Update consumer
- `js.deleteConsumer(stream, name)` - Delete consumer
- `js.getConsumerInfo(stream, name)` - Get consumer information
- `js.getConsumerNames(stream)` - List consumers for stream
- `js.pullSubscribe(stream, subject, durable)` - Create pull subscription
- `js.pullMessages(sub, batchSize, timeout)` - Pull messages
- `js.pushSubscribe(stream, subject, durable, handler)` - Create push subscription

#### Configuration
- `nats.streamConfig(options)` - Create stream configuration
- `nats.consumerConfig(options)` - Create consumer configuration
- `nats.tlsOptions(options)` - Create TLS configuration

#### Monitoring
- `js.getAccountInfo()` - Get JetStream account information

### Error Codes

The extension uses structured error codes for better debugging:
- 1001: No VU state available
- 1002: Connection is closed
- 1003: Invalid configuration
- 1004: Stream not found
- 1005: Consumer not found
- 1006: Operation timed out
- 1007: No message available
- 1008: Subject cannot be empty
- 1009: Publish failed
- 1010: Subscription failed
- 1011: Request failed
- 1012: Drain failed
- 1013: Flush failed
- 1014: JetStream not available
- 1015: Stream name cannot be empty
- 1016: Failed to add stream
- 1017: Stream not found
- 1018: Failed to update stream
- 1019: Failed to delete stream
- 1020: Failed to get stream info
- 1021: Failed to publish to JetStream
- 1022: Failed to publish async to JetStream
- 1024: Consumer name cannot be empty
- 1025: Failed to add consumer
- 1026: Consumer not found
- 1027: Failed to update consumer
- 1028: Failed to delete consumer
- 1029: Failed to get consumer info
- 1030: Failed to create pull subscription
- 1031: Subscription cannot be nil
- 1032: Failed to fetch messages
- 1033: Failed to create push subscription
- 1034: Failed to get account info
- 1035: Failed to purge stream
- 1036: Failed to delete message

## License

AGPL-3.0