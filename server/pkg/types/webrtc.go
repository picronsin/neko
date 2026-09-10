package types

import (
	"errors"

	"github.com/m1k1o/neko/server/pkg/types/codec"
)

var (
	ErrWebRTCDataChannelNotFound = errors.New("webrtc data channel not found")
	ErrWebRTCConnectionNotFound  = errors.New("webrtc connection not found")
	ErrWebRTCStreamNotFound      = errors.New("webrtc stream not found")
)

type ICEServer struct {
	URLs       []string `mapstructure:"urls"       json:"urls"`
	Username   string   `mapstructure:"username"   json:"username,omitempty"`
	Credential string   `mapstructure:"credential" json:"credential,omitempty"`
}

// SessionDescription and ICECandidate are transport-neutral WebRTC values.
// Pion conversion belongs to the WebRTC adapter, not application services.
type SessionDescription struct {
	SDP  string `json:"sdp"`
	Type string `json:"type"`
}

type ICECandidate struct {
	Candidate        string  `json:"candidate"`
	SDPMid           *string `json:"sdpMid,omitempty"`
	SDPMLineIndex    *uint16 `json:"sdpMLineIndex,omitempty"`
	UsernameFragment *string `json:"usernameFragment,omitempty"`
}

type PeerVideo struct {
	Disabled bool   `json:"disabled"`
	ID       string `json:"id"`
	Auto     bool   `json:"auto"`
}

type PeerVideoRequest struct {
	Disabled *bool           `json:"disabled,omitempty"`
	Selector *StreamSelector `json:"selector,omitempty"`
	Auto     *bool           `json:"auto,omitempty"`
}

type PeerAudio struct {
	Disabled bool `json:"disabled"`
}

type PeerAudioRequest struct {
	Disabled *bool `json:"disabled,omitempty"`
}

type WebRTCPeer interface {
	CreateOffer(ICERestart bool) (*SessionDescription, error)
	CreateAnswer() (*SessionDescription, error)
	SetRemoteDescription(SessionDescription) error
	SetCandidate(ICECandidate) error

	SetPaused(isPaused bool) error
	Paused() bool

	SetVideo(PeerVideoRequest) error
	Video() PeerVideo
	SetAudio(PeerAudioRequest) error
	Audio() PeerAudio

	SendCursorPosition(x, y int) error
	SendCursorImage(cur *CursorImage, img []byte) error

	Destroy()
}

type WebRTCManager interface {
	Start()
	Shutdown() error

	ICEServers() []ICEServer

	CreatePeer(session Session, videoCodec codec.RTPCodec) (*SessionDescription, WebRTCPeer, error)
	SetCursorPosition(x, y int)
}
