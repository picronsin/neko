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

function requiredString(payload: Record<string, unknown>, event: ProtocolEvent, field: string, allowEmpty = false) {
  const value = payload[field]
  if (typeof value !== 'string' || (!allowEmpty && value.trim() === '')) {
    throw new Error(`${event} payload field '${field}' must be a non-empty string`)
  }
}

function requiredNumber(payload: Record<string, unknown>, event: ProtocolEvent, field: string) {
  const value = payload[field]
  if (typeof value !== 'number' || !Number.isFinite(value) || value < 0) {
    throw new Error(`${event} payload field '${field}' must be a non-negative number`)
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
      for (const field of ['video', 'audio']) {
        if (value[field] !== undefined && !isRecord(value[field])) {
          throw new Error(`${event} payload field '${field}' must be an object`)
        }
      }
      return
    }
    case PROTOCOL_EVENT.SIGNAL_OFFER:
    case PROTOCOL_EVENT.SIGNAL_ANSWER: {
      const value = objectPayload(event, payload)
      rejectUnknown(value, event, ['sdp'])
      requiredString(value, event, 'sdp')
      return
    }
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
    case PROTOCOL_EVENT.SCREEN_SET: {
      const value = objectPayload(event, payload)
      rejectUnknown(value, event, ['width', 'height', 'rate'])
      requiredNumber(value, event, 'width')
      requiredNumber(value, event, 'height')
      requiredNumber(value, event, 'rate')
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
      return
    }
    case PROTOCOL_EVENT.CHAT_MESSAGE: {
      const value = objectPayload(event, payload)
      rejectUnknown(value, event, ['text'])
      requiredString(value, event, 'text')
      return
    }
    case PROTOCOL_EVENT.SYSTEM_ERROR: {
      const value = objectPayload(event, payload)
      rejectUnknown(value, event, ['title', 'code', 'message', 'retryable'])
      requiredString(value, event, 'title')
      requiredString(value, event, 'code')
      requiredString(value, event, 'message')
      if (typeof value.retryable !== 'boolean') {
        throw new Error(`${event} payload field 'retryable' must be a boolean`)
      }
      return
    }
  }
}
