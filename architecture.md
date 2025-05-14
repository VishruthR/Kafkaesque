# Broker

Represents a single queue. You can only send one request per connection for now.

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

Next up: Scan through TODOs here and throughout codebase and clean up!

TODO:

- Make message more compact for perforamnce; decide if content length field in header is actually necessary
- Implement processing multiple requests per connection (sticky connections)
- Handle multiple connections/remember connection state
- Currently using untyped queue (gross); ideally create generic implementation urself using ring buffer idea
- Support escaping `;` in PUSH_N body
- Document response types
- Write tests
- Topic deletion?
