import type { AxiosInstance } from 'axios'

import {
  Configuration,
  CurrentSessionApi,
  RoomControlApi,
  type SessionLoginRequest,
  type SessionLoginResponse,
} from '../api/generated'
import type { AuthHttpClient } from './auth'
import type { RoomHttpClient } from './room'

type RestHttpClient = AuthHttpClient &
  RoomHttpClient & {
    defaults: {
      baseURL?: string
      withCredentials?: boolean
      headers: {
        common: Record<string, unknown>
      }
    }
  }

type LoginResponse = SessionLoginResponse & {
  /** The generated allOf alias omits this optional cookie-disabled response field. */
  token?: string
}

function apiOrigin(apiURL: string) {
  return apiURL.replace(/\/api\/?$/, '')
}

/**
 * Adapts the generated OpenAPI operations to the framework-neutral REST ports.
 * The cast is intentionally isolated here: the generated client needs a full
 * AxiosInstance, while the SDK exposes only the small HTTP port it consumes.
 */
export function createGeneratedRestHttpClient(http: RestHttpClient, apiURL: string): RestHttpClient {
  const axios = http as unknown as AxiosInstance
  const basePath = apiOrigin(apiURL)
  const configuration = new Configuration({
    basePath,
    baseOptions: {
      withCredentials: true,
    },
  })
  const sessionAPI = new CurrentSessionApi(configuration, basePath, axios)
  const controlAPI = new RoomControlApi(configuration, basePath, axios)

  return {
    defaults: http.defaults,

    async post<T>(url: string, data?: unknown) {
      if (url === `${apiURL}/login`) {
        const response = await sessionAPI.login(data as SessionLoginRequest)
        return { data: response.data as T & LoginResponse }
      }
      if (url === `${apiURL}/logout`) {
        const response = await sessionAPI.logout()
        return { data: response.data as T }
      }
      if (url === `${apiURL}/room/control/request`) {
        const response = await controlAPI.controlRequest()
        return { data: response.data as T }
      }
      if (url === `${apiURL}/room/control/release`) {
        const response = await controlAPI.controlRelease()
        return { data: response.data as T }
      }
      if (url === `${apiURL}/room/control/take`) {
        const response = await controlAPI.controlTake()
        return { data: response.data as T }
      }
      if (url === `${apiURL}/room/control/reset`) {
        const response = await controlAPI.controlReset()
        return { data: response.data as T }
      }
      if (url.startsWith(`${apiURL}/room/control/give/`)) {
        const sessionID = decodeURIComponent(url.slice(`${apiURL}/room/control/give/`.length))
        const response = await controlAPI.controlGive(sessionID)
        return { data: response.data as T }
      }
      throw new Error(`Unsupported generated REST POST endpoint: ${url}`)
    },

    async get<T>(url: string) {
      if (url === `${apiURL}/room/control`) {
        const response = await controlAPI.controlStatus()
        return { data: response.data as T }
      }
      throw new Error(`Unsupported generated REST GET endpoint: ${url}`)
    },
  }
}
