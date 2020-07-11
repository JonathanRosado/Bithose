# Bithose

Bithose is a simple publish/subscribe service for building real-time applications. Servers send `messages` via
HTTP or RPC and clients subscribe to them via Websockets (with HTTP fallbacks).

### Features

- Websockets with http fallback
- Filters based on message labels
- Message delivery guarantee via configurable message buffers

### Labels

Messages may be sent with an optional set of key/value pairs called `labels`. Clients may then subscribe
to a specific set of messages based on label values. Values may be strings, numbers, or booleans.

### Examples

Send a message `POST /publish`
```
HTTP payload

{
    labels: {
        channel: 'chats',
        uid: 'SCDJCSDM',
    }
    data: "Encrypted Data"
}
```

Subscribe to a specific set of messages `WS /subscribe?filter=channel==chats&filter=uid==SCDJCSDM`
```
websocket connection

< {
    labels: {
        channel: 'chats',
        uid: 'SCDJCSDM',
    }
    data: "Encrypted Data"
  }
```

Subscribe to all messages `WS /subscribe`
```
websocket connection

< {
    labels: {
        channel: 'chats',
        uid: 'SCDJCSDM',
    }
    data: "Encrypted Data"
  }
```

### Use Cases

- Chat
- Real-time charts
- Server sent events