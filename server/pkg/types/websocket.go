package types

import (
	"net/http"

	"github.com/m1k1o/neko/server/pkg/protocol"
)

type WebSocketMessage = protocol.Envelope

type WebSocketHandler func(Session, WebSocketMessage) bool

type CheckOrigin func(r *http.Request) bool

type WebSocketPeer interface {
	Send(event string, payload any)
	Ping() error
	Destroy(reason string)
}

type WebSocketManager interface {
	Start()
	Shutdown() error
	AddHandler(handler WebSocketHandler)
	Upgrade(checkOrigin CheckOrigin) RouterHandler
}
