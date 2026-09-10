import { ProtocolErrorCode } from '../protocol/events.generated'

export interface ApiErrorPayload {
  code?: number
  message?: string
  error_code?: ProtocolErrorCode | string
}

export class ApiError extends Error {
  readonly status?: number
  readonly errorCode?: string

  constructor(message: string, status?: number, errorCode?: string) {
    super(message)
    this.name = 'ApiError'
    this.status = status
    this.errorCode = errorCode
  }
}

function isRecord(value: unknown): value is Record<string, unknown> {
  return value !== null && typeof value === 'object'
}

/** Normalize Axios/OpenAPI/fetch-like failures into the stable REST error contract. */
export function normalizeApiError(value: unknown, fallback = 'API request failed'): ApiError {
  if (value instanceof ApiError) {
    return value
  }

  const error = isRecord(value) ? value : undefined
  const response = error && isRecord(error.response) ? error.response : undefined
  const data = response && isRecord(response.data) ? response.data : undefined
  const message =
    typeof data?.message === 'string' ? data.message : typeof error?.message === 'string' ? error.message : fallback
  const status = typeof response?.status === 'number' ? response.status : undefined
  const errorCode = typeof data?.error_code === 'string' ? data.error_code : undefined
  return new ApiError(message, status, errorCode)
}
