import { normalizeApiError } from './api-error'

export interface RoomHttpClient {
  get<T>(url: string): Promise<{ data: T }>
  post<T = unknown>(url: string, data?: unknown): Promise<{ data: T }>
}

export interface ControlStatus {
  has_host: boolean
  host_id?: string
  epoch: number
}

/** REST room operations that do not depend on Vue, Vuex, or UI services. */
export class RoomClient {
  constructor(private readonly http: RoomHttpClient, private readonly apiURL: string) {}

  async controlStatus() {
    return this.request(
      () => this.http.get<ControlStatus>(`${this.apiURL}/room/control`),
      'unable to read control status',
    )
  }

  async requestControl() {
    await this.request(() => this.http.post(`${this.apiURL}/room/control/request`), 'unable to request control')
  }

  async releaseControl() {
    await this.request(() => this.http.post(`${this.apiURL}/room/control/release`), 'unable to release control')
  }

  async takeControl() {
    await this.request(() => this.http.post(`${this.apiURL}/room/control/take`), 'unable to take control')
  }

  async giveControl(sessionID: string) {
    await this.request(
      () => this.http.post(`${this.apiURL}/room/control/give/${encodeURIComponent(sessionID)}`),
      'unable to give control',
    )
  }

  async resetControl() {
    await this.request(() => this.http.post(`${this.apiURL}/room/control/reset`), 'unable to reset control')
  }

  private async request<T>(action: () => Promise<{ data: T }>, fallback: string): Promise<T> {
    try {
      const response = await action()
      return response.data
    } catch (error) {
      throw normalizeApiError(error, fallback)
    }
  }
}
