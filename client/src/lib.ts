import { accessor as neko } from './store'
import type { App, Plugin } from 'vue'

// Plugins
import Logger from './plugins/log'
import Client from './plugins/neko'
import Axios from './plugins/axios'
import Swal from './plugins/swal'
import Anime from './plugins/anime'
import { i18n } from './plugins/i18n'

// Components
import Connect from '~/components/connect.vue'
import Video from '~/components/video.vue'
import Menu from '~/components/menu.vue'
import Side from '~/components/side.vue'
import Controls from '~/components/controls.vue'
import Members from '~/components/members.vue'
import Emotes from '~/components/emotes.vue'
import About from '~/components/about.vue'
import Header from '~/components/header.vue'
import Chat from '~/components/chat.vue'
import Clipboard from '~/components/clipboard.vue'
import Emoji from '~/components/emoji.vue'
import Emote from '~/components/emote.vue'
import Context from '~/components/context.vue'
import Markdown from '~/components/markdown'
import Avatar from '~/components/avatar.vue'


// Stable SDK primitives are exported separately from the Vue components so
// embedders can own signaling, lifecycle, or input serialization themselves.
export { ConnectionStateMachine } from './sdk/connection-state'
export type { ConnectionEvent, ConnectionPhase, ConnectionTransition } from './sdk/connection-state'
export { encodeMediaInput, MEDIA_OPCODE } from './sdk/media-protocol'
export type { MediaInput, MediaInputEvent } from './sdk/media-protocol'
export { SignalingTransport } from './sdk/signaling'
export type { SignalingMessage, SignalingState, SignalingTransportOptions } from './sdk/signaling'
export { AuthClient } from './sdk/auth'
export type { AuthHttpClient, LoginResponse } from './sdk/auth'
export { ApiError, normalizeApiError } from './sdk/api-error'
export type { ApiErrorPayload } from './sdk/api-error'
export { RoomClient } from './sdk/room'
export type { RoomHttpClient, ControlStatus } from './sdk/room'
export { createGeneratedRestHttpClient } from './sdk/openapi'
export * from './api/generated'
export { PROTOCOL_ERROR, PROTOCOL_EVENT, RETRYABLE_PROTOCOL_ERRORS } from './protocol/events.generated'
export type { ProtocolErrorCode, ProtocolEvent } from './protocol/events.generated'
export * from './protocol/payloads.generated'
export { NetworkQualityMonitor, classifyNetworkQuality } from './sdk/network-monitor'
export { ControlInputController } from './sdk/control-input'
export type { ControlInputOptions, NormalizedWheel } from './sdk/control-input'
export { MediaSession } from './sdk/media-session'
export type { MediaInputData, MediaSessionOptions } from './sdk/media-session'
export { NekoClient } from './neko'
export type { NekoClientRuntime, NekoStatePort, NekoUIAdapter } from './neko/runtime'
export type {
  NetworkQuality,
  NetworkQualitySample,
  NetworkStatsPeer,
  NetworkQualityMonitorOptions,
} from './sdk/network-monitor'

/** Vue 3 plugin for embedding neko components in another application. */
export const NekoPlugin: Plugin = {
  install(app: App) {
    app.use(i18n)
    app.use(Logger)
    app.use(Axios)
    app.use(Swal)
    app.use(Anime)
    app.use(Client)
    app.config.globalProperties.$accessor = neko
  },
}

function extend(component: any) {
  return component
}

export const NekoConnect = extend(Connect)
export const NekoVideo = extend(Video)
export const NekoMenu = extend(Menu)
export const NekoSide = extend(Side)
export const NekoControls = extend(Controls)
export const NekoMembers = extend(Members)
export const NekoEmotes = extend(Emotes)
export const NekoAbout = extend(About)
export const NekoHeader = extend(Header)
export const NekoChat = extend(Chat)
export const NekoClipboard = extend(Clipboard)
export const NekoEmoji = extend(Emoji)
export const NekoEmote = extend(Emote)
export const NekoMarkdown = extend(Markdown)
export const NekoContext = extend(Context)
export const NekoAvatar = extend(Avatar)

neko.initialise()
export default neko
