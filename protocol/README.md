# Neko protocol contracts

The files in this directory are the protocol source of truth:

- `events.json` defines the WebSocket event names.
- `errors.json` defines stable error codes returned through `system/error`.
- `websocket.schema.json` defines the canonical `{event,payload}` envelope.
- `payloads.schema.json` defines the event-specific payloads for system
  lifecycle, signaling, sessions, control, screen, clipboard, keyboard,
  broadcast, opaque messages, chat, file transfer, plugins, and protocol
  errors.
- `media-input.schema.json` defines the binary WebRTC DataChannel input.

Run `go generate ./pkg/protocol` from `server/` after changing the JSON
contracts. The generator updates the Go and TypeScript protocol surfaces and
fails when an event is missing a payload contract or an explicit no-payload
declaration; do not edit generated files manually.

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

The generated Go payload structs are in
`server/pkg/protocol/payloads_generated.go`; the TypeScript payload interfaces
are in `client/src/protocol/payloads.generated.ts`. The generated REST client
is created from `server/openapi.yaml` with:

```bash
cd client
npm run generate:api
```

This command uses the official OpenAPI Generator CLI configured in
`openapitools.json`; Java is required by the CLI runtime.
