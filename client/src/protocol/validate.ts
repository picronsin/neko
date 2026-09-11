import { PROTOCOL_EVENT, ProtocolEvent } from './events.generated'

function isRecord(value: unknown): value is Record<string, unknown> {
  return value !== null && typeof value === 'object' && !Array.isArray(value)
}

function objectPayload(event: ProtocolEvent, payload: unknown): Record<string, unknown> {
  if (!isRecord(payload)) {
    throw new Error(`${event} payload must be an object`)
  }
  return payload
}

function arrayPayload(event: ProtocolEvent, payload: unknown): unknown[] {
  if (!Array.isArray(payload)) {
    throw new Error(`${event} payload must be an array`)
  }
  return payload
}

function requiredArray(payload: Record<string, unknown>, event: ProtocolEvent, field: string): unknown[] {
  const value = payload[field]
  if (!Array.isArray(value)) {
    throw new Error(`${event} payload field '${field}' must be an array`)
  }
  return value
}

function requiredString(payload: Record<string, unknown>, event: ProtocolEvent, field: string, allowEmpty = false) {
  const value = payload[field]
  if (typeof value !== 'string' || (!allowEmpty && value.trim() === '')) {
    throw new Error(`${event} payload field '${field}' must be a non-empty string`)
  }
}

function requiredNumber(payload: Record<string, unknown>, event: ProtocolEvent, field: string) {
  const value = payload[field]
  if (typeof value !== 'number' || !Number.isFinite(value) || value < 0 || !Number.isSafeInteger(value)) {
    throw new Error(`${event} payload field '${field}' must be a non-negative integer`)
  }
}

function requiredInteger(payload: Record<string, unknown>, event: ProtocolEvent, field: string) {
  const value = payload[field]
  if (typeof value !== 'number' || !Number.isFinite(value) || !Number.isSafeInteger(value)) {
    throw new Error(`${event} payload field '${field}' must be an integer`)
  }
}

function requiredPositiveNumber(payload: Record<string, unknown>, event: ProtocolEvent, field: string) {
  requiredNumber(payload, event, field)
  if ((payload[field] as number) < 1) {
    throw new Error(`${event} payload field '${field}' must be greater than zero`)
  }
}

function requiredBoolean(payload: Record<string, unknown>, event: ProtocolEvent, field: string) {
  if (typeof payload[field] !== 'boolean') {
    throw new Error(`${event} payload field '${field}' must be a boolean`)
  }
}

function validatePeerVideo(event: ProtocolEvent, payload: unknown) {
  const value = objectPayload(event, payload)
  rejectUnknown(value, event, ['disabled', 'auto', 'selector'])
  for (const field of ['disabled', 'auto']) {
    if (field in value) requiredBoolean(value, event, field)
  }
}

function validatePeerAudio(event: ProtocolEvent, payload: unknown) {
  const value = objectPayload(event, payload)
  rejectUnknown(value, event, ['disabled'])
  if ('disabled' in value) requiredBoolean(value, event, 'disabled')
}

function validatePeerVideoState(event: ProtocolEvent, payload: unknown) {
  const value = objectPayload(event, payload)
  rejectUnknown(value, event, ['disabled', 'id', 'auto'])
  requiredBoolean(value, event, 'disabled')
  requiredString(value, event, 'id', true)
  requiredBoolean(value, event, 'auto')
}

function validateScreenSize(event: ProtocolEvent, payload: unknown, withId = false) {
  const value = objectPayload(event, payload)
  rejectUnknown(value, event, withId ? ['id', 'width', 'height', 'rate'] : ['width', 'height', 'rate'])
  if (withId) requiredString(value, event, 'id')
  requiredPositiveNumber(value, event, 'width')
  requiredPositiveNumber(value, event, 'height')
  requiredPositiveNumber(value, event, 'rate')
}

function validateSettings(event: ProtocolEvent, payload: unknown, withId = false) {
  const value = objectPayload(event, payload)
  const fields = [
    'private_mode',
    'locked_logins',
    'locked_controls',
    'control_protection',
    'implicit_hosting',
    'inactive_cursors',
    'merciful_reconnect',
    'heartbeat_interval',
    'control_lease_ttl',
    'plugins',
  ]
  rejectUnknown(value, event, withId ? ['id', ...fields] : fields)
  if (withId) requiredString(value, event, 'id')
  for (const field of fields.slice(0, 7)) requiredBoolean(value, event, field)
  requiredNumber(value, event, 'heartbeat_interval')
  requiredNumber(value, event, 'control_lease_ttl')
  if (!('plugins' in value)) {
    throw new Error(`${event} payload field 'plugins' is required`)
  }
}

function validateMemberProfile(event: ProtocolEvent, payload: unknown, withId = false) {
  const value = objectPayload(event, payload)
  const fields = [
    'name',
    'avatar',
    'is_admin',
    'can_login',
    'can_connect',
    'can_watch',
    'can_host',
    'can_share_media',
    'can_access_clipboard',
    'sends_inactive_cursor',
    'can_see_inactive_cursors',
    'plugins',
  ]
  rejectUnknown(value, event, withId ? ['id', ...fields] : fields)
  if (withId) requiredString(value, event, 'id')
  requiredString(value, event, 'name', true)
  for (const field of fields.slice(2, -1)) requiredBoolean(value, event, field)
  if (!('plugins' in value)) {
    throw new Error(`${event} payload field 'plugins' is required`)
  }
  if ('avatar' in value && typeof value.avatar !== 'string') {
    throw new Error(`${event} payload field 'avatar' must be a string`)
  }
}

function validateSessionState(event: ProtocolEvent, payload: unknown, withId = false) {
  const value = objectPayload(event, payload)
  const fields = [
    'is_connected',
    'connected_since',
    'not_connected_since',
    'is_watching',
    'watching_since',
    'not_watching_since',
  ]
  rejectUnknown(value, event, withId ? ['id', ...fields] : fields)
  if (withId) requiredString(value, event, 'id')
  requiredBoolean(value, event, 'is_connected')
  requiredBoolean(value, event, 'is_watching')
  // `is_watching` is a boolean state flag, not a timestamp. Keeping it in
  // the generic timestamp loop rejected every system/init and session/state
  // payload emitted by the server.
  for (const field of ['connected_since', 'not_connected_since', 'watching_since', 'not_watching_since']) {
    if (field in value && value[field] !== null && typeof value[field] !== 'string') {
      throw new Error(`${event} payload field '${field}' must be a string or null`)
    }
  }
}

function validateSessionData(event: ProtocolEvent, payload: unknown) {
  const value = objectPayload(event, payload)
  rejectUnknown(value, event, ['id', 'profile', 'state'])
  requiredString(value, event, 'id')
  validateMemberProfile(event, value.profile)
  validateSessionState(event, value.state)
}

function validateChatHistory(event: ProtocolEvent, payload: unknown) {
  const value = objectPayload(event, payload)
  rejectUnknown(value, event, ['id', 'created', 'content', 'name', 'avatar'])
  requiredString(value, event, 'id')
  requiredString(value, event, 'created')
  const content = objectPayload(event, value.content)
  rejectUnknown(content, event, ['text'])
  requiredString(content, event, 'text', true)
  if ('name' in value && typeof value.name !== 'string') {
    throw new Error(`${event} payload field 'name' must be a string`)
  }
  if ('avatar' in value && typeof value.avatar !== 'string') {
    throw new Error(`${event} payload field 'avatar' must be a string`)
  }
}

function rejectUnknown(payload: Record<string, unknown>, event: ProtocolEvent, allowed: readonly string[]) {
  for (const key of Object.keys(payload)) {
    if (!allowed.includes(key)) {
      throw new Error(`${event} payload contains unsupported field '${key}'`)
    }
  }
}

/** Validate the event-specific payloads covered by protocol/payloads.schema.json. */
export function validateProtocolPayload(event: ProtocolEvent, payload: unknown) {
  switch (event) {
    case PROTOCOL_EVENT.CLIENT_HEARTBEAT:
    case PROTOCOL_EVENT.SYSTEM_HEARTBEAT:
    case PROTOCOL_EVENT.SIGNAL_CLOSE:
    case PROTOCOL_EVENT.CONTROL_RELEASE:
      if (payload !== undefined) {
        throw new Error(`${event} payload must be omitted`)
      }
      return
    case PROTOCOL_EVENT.SIGNAL_RESTART:
      if (payload !== undefined) {
        const value = objectPayload(event, payload)
        rejectUnknown(value, event, ['sdp'])
        requiredString(value, event, 'sdp')
      }
      return
    case PROTOCOL_EVENT.SIGNAL_REQUEST: {
      const value = objectPayload(event, payload)
      rejectUnknown(value, event, ['video_codecs', 'video', 'audio'])
      if (value.video_codecs !== undefined) {
        if (
          !Array.isArray(value.video_codecs) ||
          value.video_codecs.some((codec) => typeof codec !== 'string' || codec === '')
        ) {
          throw new Error(`${event} payload field 'video_codecs' must be a string array`)
        }
      }
      if (value.video !== undefined) validatePeerVideo(event, value.video)
      if (value.audio !== undefined) validatePeerAudio(event, value.audio)
      return
    }
    case PROTOCOL_EVENT.SIGNAL_OFFER:
    case PROTOCOL_EVENT.SIGNAL_ANSWER: {
      const value = objectPayload(event, payload)
      rejectUnknown(value, event, ['sdp'])
      requiredString(value, event, 'sdp')
      return
    }
    case PROTOCOL_EVENT.SIGNAL_VIDEO:
      // Incoming updates are the resolved PeerVideo state from the server and
      // include its stream ID; outgoing requests use PeerVideoRequest.
      validatePeerVideoState(event, payload)
      return
    case PROTOCOL_EVENT.SIGNAL_AUDIO:
      validatePeerAudio(event, payload)
      return
    case PROTOCOL_EVENT.SIGNAL_CANDIDATE: {
      const value = objectPayload(event, payload)
      rejectUnknown(value, event, ['candidate', 'sdpMid', 'sdpMLineIndex', 'usernameFragment'])
      requiredString(value, event, 'candidate')
      return
    }
    case PROTOCOL_EVENT.CONTROL_RENEW: {
      const value = objectPayload(event, payload)
      rejectUnknown(value, event, ['epoch'])
      requiredNumber(value, event, 'epoch')
      return
    }
    case PROTOCOL_EVENT.CONTROL_REQUEST: {
      if (payload === undefined) return
      const value = objectPayload(event, payload)
      rejectUnknown(value, event, ['id'])
      requiredString(value, event, 'id')
      return
    }
    case PROTOCOL_EVENT.CONTROL_HOST: {
      const value = objectPayload(event, payload)
      rejectUnknown(value, event, ['id', 'has_host', 'host_id', 'epoch'])
      requiredString(value, event, 'id')
      requiredBoolean(value, event, 'has_host')
      requiredNumber(value, event, 'epoch')
      if ('host_id' in value && typeof value.host_id !== 'string') {
        throw new Error(`${event} payload field 'host_id' must be a string`)
      }
      return
    }
    case PROTOCOL_EVENT.SCREEN_SET: {
      validateScreenSize(event, payload)
      return
    }
    case PROTOCOL_EVENT.SCREEN_UPDATED:
      validateScreenSize(event, payload, true)
      return
    case PROTOCOL_EVENT.CLIPBOARD_SET: {
      const value = objectPayload(event, payload)
      rejectUnknown(value, event, ['text'])
      requiredString(value, event, 'text', true)
      return
    }
    case PROTOCOL_EVENT.CLIPBOARD_UPDATED: {
      const value = objectPayload(event, payload)
      rejectUnknown(value, event, ['text'])
      requiredString(value, event, 'text', true)
      return
    }
    case PROTOCOL_EVENT.KEYBOARD_MAP: {
      const value = objectPayload(event, payload)
      rejectUnknown(value, event, ['layout', 'variant'])
      requiredString(value, event, 'layout')
      return
    }
    case PROTOCOL_EVENT.KEYBOARD_MODIFIERS: {
      const value = objectPayload(event, payload)
      rejectUnknown(value, event, ['shift', 'capslock', 'control', 'alt', 'numlock', 'meta', 'super', 'altgr'])
      for (const field of ['shift', 'capslock', 'control', 'alt', 'numlock', 'meta', 'super', 'altgr']) {
        if (field in value) requiredBoolean(value, event, field)
      }
      return
    }
    case PROTOCOL_EVENT.CHAT_MESSAGE: {
      const value = objectPayload(event, payload)
      if ('text' in value) {
        rejectUnknown(value, event, ['text'])
        requiredString(value, event, 'text')
        return
      }
      rejectUnknown(value, event, ['id', 'created', 'content', 'name', 'avatar'])
      requiredString(value, event, 'id')
      requiredString(value, event, 'created')
      if (!isRecord(value.content)) {
        throw new Error(`${event} payload field 'content' must be an object`)
      }
      requiredString(value.content, event, 'text', true)
      return
    }
    case PROTOCOL_EVENT.SYSTEM_LOGS: {
      const values = arrayPayload(event, payload)
      for (const item of values) {
        const value = objectPayload(event, item)
        rejectUnknown(value, event, ['level', 'fields', 'message'])
        requiredString(value, event, 'level')
        requiredString(value, event, 'message', true)
        if (value.fields !== undefined && !isRecord(value.fields)) {
          throw new Error(`${event} payload field 'fields' must be an object`)
        }
      }
      return
    }
    case PROTOCOL_EVENT.SYSTEM_ERROR: {
      const value = objectPayload(event, payload)
      rejectUnknown(value, event, ['title', 'code', 'message', 'retryable'])
      requiredString(value, event, 'title')
      requiredString(value, event, 'code')
      requiredString(value, event, 'message')
      requiredBoolean(value, event, 'retryable')
      return
    }
    case PROTOCOL_EVENT.SYSTEM_INIT: {
      const value = objectPayload(event, payload)
      rejectUnknown(value, event, [
        'session_id',
        'control_host',
        'screen_size',
        'sessions',
        'settings',
        'touch_events',
        'screencast_enabled',
        'webrtc',
      ])
      requiredString(value, event, 'session_id')
      const controlHost = objectPayload(event, value.control_host)
      rejectUnknown(controlHost, event, ['id', 'has_host', 'host_id', 'epoch'])
      requiredString(controlHost, event, 'id', true)
      requiredBoolean(controlHost, event, 'has_host')
      requiredNumber(controlHost, event, 'epoch')
      if ('host_id' in controlHost && typeof controlHost.host_id !== 'string') {
        throw new Error(`${event} payload field 'host_id' must be a string`)
      }
      validateScreenSize(event, value.screen_size)
      if (!isRecord(value.sessions)) throw new Error(`${event} payload field 'sessions' must be an object`)
      for (const session of Object.values(value.sessions)) validateSessionData(event, session)
      validateSettings(event, value.settings)
      requiredBoolean(value, event, 'touch_events')
      requiredBoolean(value, event, 'screencast_enabled')
      const webrtc = objectPayload(event, value.webrtc)
      rejectUnknown(webrtc, event, ['videos'])
      const videos = requiredArray(webrtc, event, 'videos')
      if (videos.some((video) => typeof video !== 'string')) throw new Error(`${event} payload videos must be strings`)
      return
    }
    case PROTOCOL_EVENT.SYSTEM_ADMIN: {
      const value = objectPayload(event, payload)
      rejectUnknown(value, event, [])
      return
    }
    case PROTOCOL_EVENT.SYSTEM_SETTINGS:
      validateSettings(event, payload, true)
      return
    case PROTOCOL_EVENT.SYSTEM_DISCONNECT: {
      const value = objectPayload(event, payload)
      rejectUnknown(value, event, ['message'])
      requiredString(value, event, 'message', true)
      return
    }
    case PROTOCOL_EVENT.SIGNAL_PROVIDE: {
      const value = objectPayload(event, payload)
      rejectUnknown(value, event, ['sdp', 'iceservers', 'video', 'audio'])
      requiredString(value, event, 'sdp')
      const iceServers = requiredArray(value, event, 'iceservers')
      for (const ice of iceServers) {
        const server = objectPayload(event, ice)
        rejectUnknown(server, event, ['urls', 'username', 'credential'])
        const urls = requiredArray(server, event, 'urls')
        if (urls.some((url) => typeof url !== 'string')) throw new Error(`${event} payload ICE URLs must be strings`)
      }
      validatePeerVideoState(event, value.video)
      validatePeerAudio(event, value.audio)
      return
    }
    case PROTOCOL_EVENT.SESSION_CREATED:
      validateSessionData(event, payload)
      return
    case PROTOCOL_EVENT.SESSION_DELETED:
    case PROTOCOL_EVENT.FILE_CHOOSER_DIALOG_OPENED:
    case PROTOCOL_EVENT.FILE_CHOOSER_DIALOG_CLOSED: {
      const value = objectPayload(event, payload)
      rejectUnknown(value, event, ['id'])
      requiredString(value, event, 'id', true)
      return
    }
    case PROTOCOL_EVENT.SESSION_PROFILE:
      validateMemberProfile(event, payload, true)
      return
    case PROTOCOL_EVENT.SESSION_STATE:
      validateSessionState(event, payload, true)
      return
    case PROTOCOL_EVENT.SESSION_CURSORS: {
      const values = arrayPayload(event, payload)
      for (const item of values) {
        const value = objectPayload(event, item)
        rejectUnknown(value, event, ['id', 'cursors'])
        requiredString(value, event, 'id')
        const cursors = requiredArray(value, event, 'cursors')
        for (const cursor of cursors) {
          const point = objectPayload(event, cursor)
          rejectUnknown(point, event, ['x', 'y'])
          requiredInteger(point, event, 'x')
          requiredInteger(point, event, 'y')
        }
      }
      return
    }
    case PROTOCOL_EVENT.SEND_UNICAST: {
      const value = objectPayload(event, payload)
      rejectUnknown(value, event, ['sender', 'receiver', 'subject', 'body'])
      requiredString(value, event, 'sender')
      requiredString(value, event, 'receiver')
      requiredString(value, event, 'subject', true)
      if (!('body' in value)) throw new Error(`${event} payload field 'body' is required`)
      return
    }
    case PROTOCOL_EVENT.SEND_BROADCAST: {
      const value = objectPayload(event, payload)
      rejectUnknown(value, event, ['sender', 'subject', 'body'])
      requiredString(value, event, 'sender')
      requiredString(value, event, 'subject', true)
      if (!('body' in value)) throw new Error(`${event} payload field 'body' is required`)
      return
    }
    case PROTOCOL_EVENT.CHAT_INIT: {
      const value = objectPayload(event, payload)
      rejectUnknown(value, event, ['enabled', 'history'])
      requiredBoolean(value, event, 'enabled')
      if (value.history !== undefined) {
        const history = requiredArray(value, event, 'history')
        history.forEach((message) => validateChatHistory(event, message))
      }
      return
    }
    case PROTOCOL_EVENT.CHAT_EMOTE: {
      const value = objectPayload(event, payload)
      rejectUnknown(value, event, ['id', 'emote'])
      requiredString(value, event, 'id')
      requiredString(value, event, 'emote')
      return
    }
    case PROTOCOL_EVENT.FILETRANSFER_UPDATE: {
      if (payload === undefined) return
      const value = objectPayload(event, payload)
      rejectUnknown(value, event, ['enabled', 'root_dir', 'user_download', 'user_upload', 'user_delete', 'files'])
      requiredBoolean(value, event, 'enabled')
      requiredString(value, event, 'root_dir', true)
      requiredBoolean(value, event, 'user_download')
      requiredBoolean(value, event, 'user_upload')
      requiredBoolean(value, event, 'user_delete')
      const files = requiredArray(value, event, 'files')
      for (const file of files) {
        const item = objectPayload(event, file)
        rejectUnknown(item, event, ['name', 'type', 'size'])
        requiredString(item, event, 'name', true)
        requiredString(item, event, 'type')
        if ('size' in item) requiredNumber(item, event, 'size')
      }
      return
    }
    case PROTOCOL_EVENT.OPENINAPP_INIT: {
      const value = objectPayload(event, payload)
      rejectUnknown(value, event, ['enabled'])
      requiredBoolean(value, event, 'enabled')
      return
    }
    case PROTOCOL_EVENT.OPENINAPP_OPENLINK: {
      const value = objectPayload(event, payload)
      rejectUnknown(value, event, ['text'])
      requiredString(value, event, 'text')
      return
    }
  }
}
