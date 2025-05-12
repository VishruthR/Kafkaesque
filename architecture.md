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
[...data]
</body>
```

TODO:

- Make message more compact for perforamnce
- Implement push N, pull N, pull HEAD, push, pull query support
- Implement processing multiple requests per connection (sticky connections)
