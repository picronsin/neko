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
    case PROTOCOL_EVENT.SIGNAL_RESTART:
    case PROTOCOL_EVENT.CONTROL_RELEASE:
      if (payload !== undefined) {
        throw new Error(`${event} payload must be omitted`)
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
      validatePeerVideo(event, payload)
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
    case PROTOCOL_EVENT.SCREEN_SET: {
      const value = objectPayload(event, payload)
      rejectUnknown(value, event, ['width', 'height', 'rate'])
      requiredPositiveNumber(value, event, 'width')
      requiredPositiveNumber(value, event, 'height')
      requiredPositiveNumber(value, event, 'rate')
      return
    }
    case PROTOCOL_EVENT.CLIPBOARD_SET: {
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
    case PROTOCOL_EVENT.SYSTEM_INIT:
    case PROTOCOL_EVENT.SYSTEM_ADMIN:
    case PROTOCOL_EVENT.SYSTEM_SETTINGS:
    case PROTOCOL_EVENT.SYSTEM_DISCONNECT:
    case PROTOCOL_EVENT.SIGNAL_PROVIDE:
    case PROTOCOL_EVENT.SESSION_CREATED:
    case PROTOCOL_EVENT.SESSION_DELETED:
    case PROTOCOL_EVENT.SESSION_PROFILE:
    case PROTOCOL_EVENT.SESSION_STATE:
    case PROTOCOL_EVENT.CONTROL_HOST:
    case PROTOCOL_EVENT.SCREEN_UPDATED:
    case PROTOCOL_EVENT.CLIPBOARD_UPDATED:
    case PROTOCOL_EVENT.BROADCAST_STATUS:
    case PROTOCOL_EVENT.SEND_UNICAST:
    case PROTOCOL_EVENT.SEND_BROADCAST:
    case PROTOCOL_EVENT.FILE_CHOOSER_DIALOG_OPENED:
    case PROTOCOL_EVENT.FILE_CHOOSER_DIALOG_CLOSED:
    case PROTOCOL_EVENT.CHAT_INIT:
    case PROTOCOL_EVENT.CHAT_EMOTE:
    case PROTOCOL_EVENT.FILETRANSFER_UPDATE:
    case PROTOCOL_EVENT.OPENINAPP_INIT:
    case PROTOCOL_EVENT.OPENINAPP_OPENLINK:
      objectPayload(event, payload)
      return
    case PROTOCOL_EVENT.SESSION_CURSORS:
      arrayPayload(event, payload)
      return
  }
}
