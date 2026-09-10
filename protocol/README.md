# Neko protocol contracts

The files in this directory are the protocol source of truth:

- `events.json` defines the WebSocket event names.
- `errors.json` defines stable error codes returned through `system/error`.
- `websocket.schema.json` defines the canonical `{event,payload}` envelope.
- `media-input.schema.json` defines the binary WebRTC DataChannel input.

Run `go generate ./pkg/protocol` from `server/` after changing the JSON
contracts. The generator updates the Go and TypeScript protocol surfaces; do
not edit generated files manually.

`media-input.schema.json` is the logical contract for control input shared by
the browser and server. The WebRTC DataChannel representation is a network
byte-order packet:

```text
opcode:uint8 | length:uint16 | epoch:uint64 | payload
```

`length` counts the epoch and payload bytes after the three-byte header. The
server accepts desktop input only when the packet epoch matches the current
control lease. Inactive cursor movement remains available to viewers without a
lease; all other input is rejected unless the sender is the current holder.

The control lease is exposed through the REST control status response and the
`system/init` and `control/host` WebSocket payloads as `epoch`.

While a session holds control, the client periodically sends
`control/renew` with the current epoch. The server renews the lease only for
the current holder; an old session or stale epoch cannot keep ownership alive.
