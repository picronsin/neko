import { NekoClientRuntime } from './runtime'

type Translator = (key: string, params?: Record<string, unknown>) => string

/** Adapts Vue 3 global properties without leaking UI details into the SDK. */
export function createVueNekoRuntime(vue: any, fallbackTranslate?: Translator): NekoClientRuntime {
  // Vue 3's `app.config.globalProperties` does not expose vue-i18n's `$t` as
  // a callable function in every setup.  Prefer it when available, then use
  // the i18n instance supplied by the application, and finally fall back to
  // the key so a notification can never crash the UI.
  const translate: Translator =
    typeof vue.$t === 'function'
      ? vue.$t.bind(vue)
      : fallbackTranslate || ((key: string) => key)

  return {
    http: vue.$http,
    state: vue.$accessor,
    ui: {
      translate(key, params) {
        return translate(key, params)
      },
      notify(options) {
        vue.$notify(options)
      },
      alert(options) {
        vue.$swal(options)
      },
      log: vue.$log,
    },
  }
}
