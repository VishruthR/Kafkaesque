# Broker

The broker hosts all topics for a given partition. Publishers push messages to a topic and consumers pull messages from a topic.

### Topics

Each topic is essentially a queue of messages. Publishers publish to topics and consumers consume from topics. If a message is consumed from a topic, it is no longer available in that topic. All messages in a topic are durable.

### Message format

```
<header>
PUSH/PULL
[topic]
[content length in bytes]
</header>
<body>
[...data] // If you are doing a PUSH_N, each element in the body should be delimited by a `;`
</body>
```

Next up: Make topics durable: Use WAL w/ checkpointing. Then write tests for current system

TODO:

- Make message more compact for perforamnce; decide if content length field in header is actually necessary
- Implement processing multiple requests per connection (sticky connections)
- Handle multiple connections/remember connection state
- Currently using untyped queue (gross); ideally create generic implementation urself using ring buffer idea
- Support escaping `;` in PUSH_N body
- Document response types
- Write tests
- Topic deletion?
- Implement RAFT for consensus (consider minikube or alternatives, find a way to run distributed system experiments on your local machine; Might require getting a new computer)
- Support for pushing to multiple topics at the same time (potentially wildcard matching topics?)

RAFT: Fault-tolerance/availability + Partioning to distrbute load!
