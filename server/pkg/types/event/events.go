// Package event keeps the historical event names used by the server while
// sourcing their values from the generated protocol catalog.
package event

import "github.com/m1k1o/neko/server/pkg/protocol"

const (
	SYSTEM_INIT       = protocol.SystemInit
	SYSTEM_ADMIN      = protocol.SystemAdmin
	SYSTEM_SETTINGS   = protocol.SystemSettings
	SYSTEM_LOGS       = protocol.SystemLogs
	SYSTEM_DISCONNECT = protocol.SystemDisconnect
	SYSTEM_HEARTBEAT  = protocol.SystemHeartbeat
	SYSTEM_ERROR      = protocol.SystemError
)

const CLIENT_HEARTBEAT = protocol.ClientHeartbeat

const (
	SIGNAL_REQUEST   = protocol.SignalRequest
	SIGNAL_RESTART   = protocol.SignalRestart
	SIGNAL_OFFER     = protocol.SignalOffer
	SIGNAL_ANSWER    = protocol.SignalAnswer
	SIGNAL_PROVIDE   = protocol.SignalProvide
	SIGNAL_CANDIDATE = protocol.SignalCandidate
	SIGNAL_CLOSE     = protocol.SignalClose
	SIGNAL_VIDEO     = protocol.SignalVideo
	SIGNAL_AUDIO     = protocol.SignalAudio
)

const (
	SESSION_CREATED = protocol.SessionCreated
	SESSION_DELETED = protocol.SessionDeleted
	SESSION_PROFILE = protocol.SessionProfile
	SESSION_STATE   = protocol.SessionState
	SESSION_CURSORS = protocol.SessionCursors
)

const (
	CONTROL_HOST    = protocol.ControlHost
	CONTROL_RELEASE = protocol.ControlRelease
	CONTROL_REQUEST = protocol.ControlRequest
	CONTROL_RENEW   = protocol.ControlRenew
)

const (
	SCREEN_UPDATED = protocol.ScreenUpdated
	SCREEN_SET     = protocol.ScreenSet
)

const (
	CLIPBOARD_UPDATED = protocol.ClipboardUpdated
	CLIPBOARD_SET     = protocol.ClipboardSet
)

const (
	KEYBOARD_MODIFIERS = protocol.KeyboardModifiers
	KEYBOARD_MAP       = protocol.KeyboardMap
)

const (
	BROADCAST_STATUS = protocol.BroadcastStatus
	SEND_UNICAST     = protocol.SendUnicast
	SEND_BROADCAST   = protocol.SendBroadcast
)

const (
	FILE_CHOOSER_DIALOG_OPENED = protocol.FileChooserDialogOpened
	FILE_CHOOSER_DIALOG_CLOSED = protocol.FileChooserDialogClosed
)
