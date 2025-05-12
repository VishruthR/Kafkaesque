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

TODO:

- Make message more compact for perforamnce; decide if content length field in header is actually necessary
- Implement push N, pull N, pull HEAD, push, pull query support
- Implement processing multiple requests per connection (sticky connections)
- Handle multiple connections/remember connection state
- Support topic creation
- Currently using untyped queue (gross); ideally create generic implementation urself using ring buffer idea
- Support delimiting `;` in PUSH_N body
- Document response types
- Write tests
