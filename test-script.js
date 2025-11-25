import nats from "k6/Pondigo/nats";

export const options = {
  vus: 1,
  duration: "30s",
};

export default function () {
  // Test connection
  const conn = nats.connect({
    urls: ["nats://localhost:4222"],
    reconnectWait: 2,
    maxReconnects: 10,
    pingInterval: 60,
    maxPingsOut: 2,
  });

  if (!conn) {
    console.log("Failed to connect to NATS");
    return;
  }

  console.log("Connected to NATS successfully");

  // Test basic publish
  const publishResult = conn.publish("test.subject", "Hello from k6!");
  if (publishResult) {
    console.log("Published message successfully");
  } else {
    console.log("Failed to publish message");
  }

  // Test JetStream
  const js = nats.jetStream(conn);
  if (!js) {
    console.log("Failed to create JetStream context");
    conn.close();
    return;
  }

  console.log("JetStream context created successfully");

  // Create stream config
  const streamConfig = nats.streamConfig({
    name: "TEST_STREAM",
    subjects: ["test.>"],
    retention: "limits",
    storage: "file",
    replicas: 1,
    maxBytes: 1048576, // 1MB
    maxMsgs: 1000,
    maxAge: 3600, // 1 hour
    discard: "old",
  });

  // Add stream
  try {
    js.addStream(streamConfig);
    console.log("Stream created successfully");
  } catch (error) {
    console.log("Failed to create stream:", error.message);
  }

  // Publish to JetStream
  try {
    js.publish("test.js", "Hello JetStream!");
    console.log("Published to JetStream successfully");
  } catch (error) {
    console.log("Failed to publish to JetStream:", error.message);
  }

  // Get stream info
  try {
    const streamInfo = js.getStreamInfo("TEST_STREAM");
    if (streamInfo) {
      console.log(
        "Stream info:",
        streamInfo.config.name,
        "Messages:",
        streamInfo.state.messages,
      );
    }
  } catch (error) {
    console.log("Failed to get stream info:", error.message);
  }

  // Create consumer config
  const consumerConfig = nats.consumerConfig({
    stream: "TEST_STREAM",
    durable: "TEST_CONSUMER",
    deliverPolicy: "all",
    ackPolicy: "explicit",
    ackWait: 30,
    maxDeliver: 3,
    replayPolicy: "instant",
  });

  // Add consumer
  try {
    js.addConsumer("TEST_STREAM", consumerConfig);
    console.log("Consumer created successfully");
  } catch (error) {
    console.log("Failed to create consumer:", error.message);
  }

  // Get consumer info
  try {
    const consumerInfo = js.getConsumerInfo("TEST_STREAM", "TEST_CONSUMER");
    if (consumerInfo) {
      console.log(
        "Consumer info:",
        consumerInfo.config.durable,
        "Delivered:",
        consumerInfo.delivered,
      );
    }
  } catch (error) {
    console.log("Failed to get consumer info:", error.message);
  }

  // Test subscription
  let messageReceived = false;
  const subscription = conn.subscribe("test.subscription", "", (msg) => {
    console.log("Received message:", msg.subject, String(msg.data));
    messageReceived = true;
  });

  if (subscription) {
    console.log("Subscription created successfully");

    // Publish a message to test subscription
    conn.publish("test.subscription", "Test subscription message");

    // Wait a bit for message to be received
    setTimeout(() => {
      if (messageReceived) {
        console.log("Subscription test passed");
      } else {
        console.log("Subscription test failed - no message received");
      }
    }, 100);
  }

  // Test request/reply
  try {
    const reply = conn.request("test.request", "ping", 5000);
    if (reply) {
      console.log("Request/reply test successful:", String(reply.data));
    } else {
      console.log("Request/reply test failed - no reply");
    }
  } catch (error) {
    console.log("Request/reply test failed:", error.message);
  }

  // Clean up
  try {
    js.deleteConsumer("TEST_STREAM", "TEST_CONSUMER");
    console.log("Consumer deleted successfully");
  } catch (error) {
    console.log("Failed to delete consumer:", error.message);
  }

  try {
    js.deleteStream("TEST_STREAM");
    console.log("Stream deleted successfully");
  } catch (error) {
    console.log("Failed to delete stream:", error.message);
  }

  conn.close();
  console.log("Connection closed");
}