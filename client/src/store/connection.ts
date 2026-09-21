import { getterTree, mutationTree, actionTree } from './helpers'
import { NetworkPath, NetworkProtocol } from '~/sdk/network-monitor'

export const namespaced = true

export type ConnectionState = 'disconnected' | 'connecting' | 'reconnecting' | 'connected'
export type NetworkQuality = 'unknown' | 'good' | 'fair' | 'poor'

export const state = () => ({
  state: 'disconnected' as ConnectionState,
  // Authentication completes as soon as the authenticated websocket sends
  // system/init. Media transport may still be negotiating at that point.
  authenticated: false,
  // A reconnecting ICE transport still owns a live session. Keep this
  // separate from the lifecycle label so the UI does not tear down the
  // session during a transient network interruption.
  connected: false,
  error: '',
  quality: 'unknown' as NetworkQuality,
  rtt: null as number | null,
  path: 'unknown' as NetworkPath,
  protocol: 'unknown' as NetworkProtocol,
})

export const getters = getterTree(state, {
  connected: (state) => state.connected,
  authenticated: (state) => state.authenticated,
  connecting: (state) => state.state === 'connecting' || state.state === 'reconnecting',
})

export const mutations = mutationTree(state, {
  setConnecting(state) {
    state.state = 'connecting'
    state.authenticated = false
    state.connected = false
    state.error = ''
    state.quality = 'unknown'
    state.rtt = null
    state.path = 'unknown'
    state.protocol = 'unknown'
  },

  setConnected(state, connected: boolean) {
    state.state = connected ? 'connected' : 'disconnected'
    state.authenticated = connected
    state.connected = connected
    if (!connected) {
      state.quality = 'unknown'
      state.rtt = null
      state.path = 'unknown'
      state.protocol = 'unknown'
    }
  },

  setState(state, connectionState: ConnectionState) {
    state.state = connectionState
    if (connectionState === 'connected') {
      state.connected = true
    } else if (connectionState === 'connecting' || connectionState === 'disconnected') {
      state.connected = false
    }
    if (connectionState === 'disconnected') {
      state.quality = 'unknown'
      state.rtt = null
      state.path = 'unknown'
      state.protocol = 'unknown'
    }
  },

  setAuthenticated(state, authenticated: boolean) {
    state.authenticated = authenticated
    if (!authenticated) {
      state.connected = false
    }
  },

  setError(state, message: string) {
    state.error = message
  },

  setNetworkQuality(state, { quality, rtt }: { quality: NetworkQuality; rtt: number | null }) {
    state.quality = quality
    state.rtt = rtt
  },

  setNetworkPath(state, { path, protocol }: { path: NetworkPath; protocol: NetworkProtocol }) {
    state.path = path
    state.protocol = protocol
  },
})

export const actions = actionTree({ state, getters, mutations }, {})
