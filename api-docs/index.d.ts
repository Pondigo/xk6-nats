/**
 * @packageDocumentation
 * xk6-nats is a k6 extension to load test NATS
 */

/**
 * @module k6/Pondigo/nats
 * @description
 * The xk6-nats project is a k6 extension that enables k6 users to load test NATS using connections, publishers, subscribers, and JetStream functionality.
 * This documentation refers to the development version of xk6-nats project, which means the latest changes and might not be released yet.
 * @see {@link https://github.com/pondigo/xk6-nats}
 */

/* TLS versions for creating a secure communication channel with NATS. */
export enum TLS_VERSIONS {
  TLS_1_0 = "tlsv1.0",
  TLS_1_1 = "tlsv1.1",
  TLS_1_2 = "tlsv1.2",
  TLS_1_3 = "tlsv1.3",
}

/* Reconnect policies for handling connection failures. */
export enum RECONNECT_POLICIES {
  RECONNECT_POLICY_NONE = "reconnect_policy_none",
  RECONNECT_POLICY_RECONNECT = "reconnect_policy_reconnect",
  RECONNECT_POLICY_RECONNECT_FOREVER = "reconnect_policy_reconnect_forever",
}

/* JetStream storage types for stream persistence. */
export enum STORAGE_TYPES {
  STORAGE_TYPE_FILE = "file",
  STORAGE_TYPE_MEMORY = "memory",
}

/* JetStream retention policies for streams. */
export enum RETENTION_POLICIES {
  RETENTION_POLICY_LIMITS = "limits",
  RETENTION_POLICY_INTEREST = "interest",
  RETENTION_POLICY_WORK_QUEUE = "workqueue",
}

/* JetStream discard policies for streams. */
export enum DISCARD_POLICIES {
  DISCARD_POLICY_OLD = "old",
  DISCARD_POLICY_NEW = "new",
}

/* JetStream deliver policies for consumers. */
export enum DELIVER_POLICIES {
  DELIVER_POLICY_ALL = "all",
  DELIVER_POLICY_LAST = "last",
  DELIVER_POLICY_NEW = "new",
  DELIVER_POLICY_BY_START_SEQUENCE = "by_start_sequence",
  DELIVER_POLICY_BY_START_TIME = "by_start_time",
}

/* JetStream acknowledgment policies for consumers. */
export enum ACK_POLICIES {
  ACK_POLICY_NONE = "none",
  ACK_POLICY_ALL = "all",
  ACK_POLICY_EXPLICIT = "explicit",
}

/* JetStream replay policies for consumers. */
export enum REPLAY_POLICIES {
  REPLAY_POLICY_INSTANT = "instant",
  REPLAY_POLICY_ORIGINAL = "original",
}

/* SASL mechanisms for authenticating to NATS. */
export enum SASL_MECHANISMS {
  SASL_NONE = "none",
  SASL_PLAIN = "plain",
  SASL_SCRAM_SHA256 = "scram_sha256",
  SASL_SCRAM_SHA512 = "scram_sha512",
  SASL_TOKEN = "token",
}

/* Time units for use in timeouts and intervals. */
export enum TIME {
  NANOSECOND = 1,
  MICROSECOND = 1000,
  MILLISECOND = 1000000,
  SECOND = 1000000000,
  MINUTE = 60000000000,
  HOUR = 3600000000000,
}

/* TLS configurations for creating a secure communication channel with NATS. */
export interface TLSConfig {
  /** Enable TLS encryption */
  enabled: boolean;
  /** Skip TLS certificate verification (insecure) */
  insecureSkipVerify: boolean;
  /** Minimum TLS version */
  minVersion: TLS_VERSIONS;
  /** Path to client certificate file */
  certFile: string;
  /** Path to client key file */
  keyFile: string;
  /** Path to CA certificate file */
  caFile: string;
}

/* SASL configurations for authenticating to NATS. */
export interface SASLConfig {
  /** SASL mechanism to use */
  mechanism: SASL_MECHANISMS;
  /** Username for authentication */
  username: string;
  /** Password for authentication */
  password: string;
  /** Token for authentication */
  token: string;
}

/* Connection configurations for connecting to NATS servers. */
export interface ConnectionConfig {
  /** List of NATS server URLs */
  urls: string[];
  /** Maximum number of reconnect attempts */
  maxReconnects: number;
  /** Time to wait between reconnect attempts in seconds */
  reconnectWait: number;
  /** Ping interval in seconds */
  pingInterval: number;
  /** Maximum outstanding pings */
  maxPingsOut: number;
  /** Allow reconnection */
  allowReconnect: boolean;
  /** TLS configuration */
  tls: TLSConfig;
  /** SASL configuration */
  sasl: SASLConfig;
  /** Connection name for identification */
  name: string;
}

/* Message format for NATS messages. */
export interface Message {
  /** Subject the message was published to */
  subject: string;
  /** Message payload data */
  data: Uint8Array;
  /** Reply subject for request/reply pattern */
  reply: string;
  /** Message headers */
  headers: Map<string, string>;
  /** Message timestamp */
  timestamp: Date;
}

/* Configuration for publishing messages. */
export interface PublishConfig {
  /** Subject to publish to */
  subject: string;
  /** Message payload data */
  data: Uint8Array;
  /** Reply subject for request/reply */
  reply: string;
  /** Message headers */
  headers: Map<string, string>;
}

/* Configuration for subscribing to messages. */
export interface SubscribeConfig {
  /** Subject pattern to subscribe to */
  subject: string;
  /** Queue group name for load balancing */
  queue: string;
  /** Message handler function */
  handler: (msg: Message) => void;
}

/* Configuration for request/reply pattern. */
export interface RequestConfig {
  /** Subject to send request to */
  subject: string;
  /** Request payload data */
  data: Uint8Array;
  /** Timeout in milliseconds */
  timeout: number;
  /** Request headers */
  headers: Map<string, string>;
}

/* JetStream stream configuration. */
export interface StreamConfig {
  /** Stream name */
  name: string;
  /** List of subjects for the stream */
  subjects: string[];
  /** Retention policy */
  retention: RETENTION_POLICIES;
  /** Maximum bytes in stream */
  maxBytes: number;
  /** Maximum messages in stream */
  maxMsgs: number;
  /** Maximum age of messages in seconds */
  maxAge: number;
  /** Number of replicas for HA */
  replicas: number;
  /** Discard policy */
  discard: DISCARD_POLICIES;
  /** Storage type */
  storage: STORAGE_TYPES;
}

/* JetStream consumer configuration. */
export interface ConsumerConfig {
  /** Stream name */
  stream: string;
  /** Consumer name */
  name: string;
  /** Durable consumer name */
  durable: string;
  /** Deliver policy */
  deliverPolicy: DELIVER_POLICIES;
  /** Start sequence number */
  optStartSeq: number;
  /** Start time */
  optStartTime: number;
  /** Acknowledgment policy */
  ackPolicy: ACK_POLICIES;
  /** Acknowledgment wait time in seconds */
  ackWait: number;
  /** Maximum delivery attempts */
  maxDeliver: number;
  /** Backoff intervals in seconds */
  backOff: number[];
  /** Subject filter */
  filterSubject: string;
  /** Replay policy */
  replayPolicy: REPLAY_POLICIES;
  /** Sample frequency */
  sampleFreq: string;
}

/* Configuration for pull consumer. */
export interface PullConfig {
  /** Batch size to pull */
  batchSize: number;
  /** Timeout in milliseconds */
  timeout: number;
  /** Maximum wait time in milliseconds */
  maxWait: number;
}

/* Configuration for push consumer. */
export interface PushConfig {
  /** Subject to subscribe to */
  subject: string;
  /** Queue group name */
  queue: string;
  /** Durable consumer name */
  durable: string;
  /** Message handler function */
  handler: (msg: Message) => void;
}

/**
 * @class
 * @classdesc Connection represents a connection to NATS servers.
 * @example
 *
 * ```javascript
 * // In init context
 * const connection = new Connection({
 *   urls: ["nats://localhost:4222"],
 *   maxReconnects: 10,
 *   reconnectWait: 2,
 * });
 *
 * // In VU code (default function)
 * connection.publish({
 *   subject: "test.subject",
 *   data: new TextEncoder().encode("Hello NATS!")
 * });
 *
 * // In teardown function
 * connection.close();
 * ```
 */
export class Connection {
  /**
   * @constructor
   * Create a new Connection.
   * @param {ConnectionConfig} connectionConfig - Connection configuration.
   * @returns {Connection} - Connection instance.
   */
  constructor(connectionConfig: ConnectionConfig);

  /**
   * @method
   * Publish a message to a subject.
   * @param {PublishConfig} publishConfig - Publish configuration.
   * @returns {void} - Nothing.
   */
  publish(publishConfig: PublishConfig): void;

  /**
   * @method
   * Subscribe to a subject pattern.
   * @param {SubscribeConfig} subscribeConfig - Subscribe configuration.
   * @returns {Subscription} - Subscription instance.
   */
  subscribe(subscribeConfig: SubscribeConfig): Subscription;

  /**
   * @method
   * Send a request and wait for a reply.
   * @param {RequestConfig} requestConfig - Request configuration.
   * @returns {Message} - Reply message.
   */
  request(requestConfig: RequestConfig): Message;

  /**
   * @method
   * Get JetStream context for stream operations.
   * @returns {JetStream} - JetStream instance.
   */
  jetStream(): JetStream;

  /**
   * @method
   * Check if connection is active.
   * @returns {boolean} - Connection status.
   */
  isConnected(): boolean;

  /**
   * @method
   * Get connection statistics.
   * @returns {ConnectionStats} - Connection statistics.
   */
  stats(): ConnectionStats;

  /**
   * @destructor
   * @description Close the connection.
   * @returns {void} - Nothing.
   */
  close(): void;
}

/**
 * @class
 * @classdesc Subscription represents a subscription to a NATS subject.
 * @example
 *
 * ```javascript
 * const subscription = connection.subscribe({
 *   subject: "test.>",
 *   handler: (msg) => {
 *     console.log(`Received: ${new TextDecoder().decode(msg.data)}`);
 *   }
 * });
 *
 * // Later...
 * subscription.unsubscribe();
 * ```
 */
export class Subscription {
  /**
   * @method
   * Unsubscribe from the subject.
   * @returns {void} - Nothing.
   */
  unsubscribe(): void;

  /**
   * @method
   * Get subscription statistics.
   * @returns {SubscriptionStats} - Subscription statistics.
   */
  stats(): SubscriptionStats;
}

/**
 * @class
 * @classdesc JetStream provides access to NATS JetStream functionality.
 * @example
 *
 * ```javascript
 * const js = connection.jetStream();
 *
 * // Create a stream
 * js.addStream({
 *   name: "TEST_STREAM",
 *   subjects: ["test.>"],
 *   storage: STORAGE_TYPES.STORAGE_TYPE_FILE,
 *   retention: RETENTION_POLICIES.RETENTION_POLICY_LIMITS
 * });
 *
 * // Publish to stream
 * js.publish("TEST_STREAM", "test.subject", new TextEncoder().encode("Hello JetStream!"));
 * ```
 */
export class JetStream {
  /**
   * @method
   * Add a new stream.
   * @param {StreamConfig} streamConfig - Stream configuration.
   * @returns {StreamInfo} - Stream information.
   */
  addStream(streamConfig: StreamConfig): StreamInfo;

  /**
   * @method
   * Delete a stream.
   * @param {string} streamName - Stream name.
   * @returns {void} - Nothing.
   */
  deleteStream(streamName: string): void;

  /**
   * @method
   * Get stream information.
   * @param {string} streamName - Stream name.
   * @returns {StreamInfo} - Stream information.
   */
  streamInfo(streamName: string): StreamInfo;

  /**
   * @method
   * Add a consumer to a stream.
   * @param {string} streamName - Stream name.
   * @param {ConsumerConfig} consumerConfig - Consumer configuration.
   * @returns {ConsumerInfo} - Consumer information.
   */
  addConsumer(streamName: string, consumerConfig: ConsumerConfig): ConsumerInfo;

  /**
   * @method
   * Delete a consumer from a stream.
   * @param {string} streamName - Stream name.
   * @param {string} consumerName - Consumer name.
   * @returns {void} - Nothing.
   */
  deleteConsumer(streamName: string, consumerName: string): void;

  /**
   * @method
   * Publish a message to a stream.
   * @param {string} streamName - Stream name.
   * @param {string} subject - Subject.
   * @param {Uint8Array} data - Message data.
   * @returns {PublishAck} - Publish acknowledgment.
   */
  publish(streamName: string, subject: string, data: Uint8Array): PublishAck;

  /**
   * @method
   * Create a pull consumer.
   * @param {string} streamName - Stream name.
   * @param {string} subject - Subject.
   * @param {string} durable - Durable name.
   * @returns {PullConsumer} - Pull consumer instance.
   */
  pullSubscribe(
    streamName: string,
    subject: string,
    durable: string,
  ): PullConsumer;

  /**
   * @method
   * Create a push consumer.
   * @param {PushConfig} pushConfig - Push configuration.
   * @returns {PushConsumer} - Push consumer instance.
   */
  pushSubscribe(pushConfig: PushConfig): PushConsumer;
}

/**
 * @class
 * @classdesc PullConsumer for consuming messages from JetStream using pull mode.
 * @example
 *
 * ```javascript
 * const consumer = js.pullSubscribe("TEST_STREAM", "test.subject", "my-durable");
 *
 * const messages = consumer.pull({
 *   batchSize: 10,
 *   timeout: 5000
 * });
 *
 * for (const msg of messages) {
 *   console.log(new TextDecoder().decode(msg.data));
 *   msg.ack();
 * }
 * ```
 */
export class PullConsumer {
  /**
   * @method
   * Pull messages from the consumer.
   * @param {PullConfig} pullConfig - Pull configuration.
   * @returns {Message[]} - Array of messages.
   */
  pull(pullConfig: PullConfig): Message[];

  /**
   * @method
   * Acknowledge a message.
   * @returns {void} - Nothing.
   */
  ack(): void;

  /**
   * @method
   * Negative acknowledge a message.
   * @returns {void} - Nothing.
   */
  nack(): void;

  /**
   * @method
   * Get consumer information.
   * @returns {ConsumerInfo} - Consumer information.
   */
  info(): ConsumerInfo;
}

/**
 * @class
 * @classdesc PushConsumer for consuming messages from JetStream using push mode.
 * @example
 *
 * ```javascript
 * const consumer = js.pushSubscribe({
 *   streamName: "TEST_STREAM",
 *   subject: "test.subject",
 *   durable: "my-durable",
 *   handler: (msg) => {
 *     console.log(new TextDecoder().decode(msg.data));
 *     msg.ack();
 *   }
 * });
 * ```
 */
export class PushConsumer {
  /**
   * @method
   * Acknowledge a message.
   * @returns {void} - Nothing.
   */
  ack(): void;

  /**
   * @method
   * Negative acknowledge a message.
   * @returns {void} - Nothing.
   */
  nack(): void;

  /**
   * @method
   * Stop the consumer.
   * @returns {void} - Nothing.
   */
  stop(): void;

  /**
   * @method
   * Get consumer information.
   * @returns {ConsumerInfo} - Consumer information.
   */
  info(): ConsumerInfo;
}

/* Connection statistics. */
export interface ConnectionStats {
  /** Number of messages sent */
  messagesSent: number;
  /** Number of bytes sent */
  bytesSent: number;
  /** Number of messages received */
  messagesReceived: number;
  /** Number of bytes received */
  bytesReceived: number;
  /** Number of reconnects */
  reconnects: number;
  /** Number of errors */
  errors: number;
  /** Connection uptime in seconds */
  uptime: number;
}

/* Subscription statistics. */
export interface SubscriptionStats {
  /** Number of messages delivered */
  messagesDelivered: number;
  /** Number of bytes delivered */
  bytesDelivered: number;
  /** Number of pending messages */
  pendingMessages: number;
  /** Number of pending bytes */
  pendingBytes: number;
  /** Subscription active status */
  active: boolean;
}

/* Stream information. */
export interface StreamInfo {
  /** Stream name */
  name: string;
  /** Stream subjects */
  subjects: string[];
  /** Stream configuration */
  config: StreamConfig;
  /** Stream state */
  state: StreamState;
  /** Creation time */
  created: Date;
}

/* Stream state information. */
export interface StreamState {
  /** Number of messages */
  messages: number;
  /** Number of bytes */
  bytes: number;
  /** First sequence */
  firstSeq: number;
  /** Last sequence */
  lastSeq: number;
  /** Consumer count */
  consumers: number;
  /** Last activity time */
  lastActivity: Date;
}

/* Consumer information. */
export interface ConsumerInfo {
  /** Stream name */
  streamName: string;
  /** Consumer name */
  name: string;
  /** Consumer configuration */
  config: ConsumerConfig;
  /** Consumer state */
  state: ConsumerState;
  /** Creation time */
  created: Date;
}

/* Consumer state information. */
export interface ConsumerState {
  /** Number of delivered messages */
  delivered: number;
  /** Number of pending messages */
  pending: number;
  /** Number of redelivered messages */
  redelivered: number;
  /** Last activity time */
  lastActivity: Date;
}

/* Publish acknowledgment. */
export interface PublishAck {
  /** Stream name */
  stream: string;
  /** Sequence number */
  seq: number;
  /** Duplicate flag */
  duplicate: boolean;
  /** Domain name */
  domain: string;
}
