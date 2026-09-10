import { normalizeApiError } from './api-error'

export interface AuthHttpClient {
  defaults: {
    headers: {
      common: Record<string, unknown>
    }
  }
  post<T>(url: string, data?: unknown): Promise<{ data: T }>
}

export interface LoginResponse {
  token?: string
}

/** Owns session authentication without depending on Vue, Vuex, or dialogs. */
export class AuthClient {
  private _token = ''

  constructor(private readonly http: AuthHttpClient, private readonly apiURL: string) {}

  get token() {
    return this._token
  }

  async login(username: string, password: string) {
    try {
      const response = await this.http.post<LoginResponse>(`${this.apiURL}/login`, {
        username,
        password,
      })
      this._token = response.data.token || ''
      this.applyToken()
      return this._token
    } catch (error) {
      throw normalizeApiError(error, 'login failed')
    }
  }

  async logout() {
    try {
      await this.http.post(`${this.apiURL}/logout`)
    } catch (error) {
      throw normalizeApiError(error, 'logout failed')
    } finally {
      this.clear()
    }
  }

  clear() {
    this._token = ''
    this.applyToken()
  }

  private applyToken() {
    const headers = this.http.defaults.headers.common
    if (this._token) {
      headers.Authorization = `Bearer ${this._token}`
    } else {
      delete headers.Authorization
    }
  }
}
