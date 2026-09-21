import type { Plugin } from 'vue'
import { NekoClient } from '~/neko'

declare global {
  const $client: NekoClient

  interface Window {
    $client: any
  }
}

declare module '@vue/runtime-core' {
  interface ComponentCustomProperties {
    $client: NekoClient
  }
}

const plugin: Plugin = {
  install(app) {
    window.$client = new NekoClient()
      .on('error', window.$log.error)
      .on('warn', window.$log.warn)
      .on('info', window.$log.info)
      .on('debug', window.$log.debug)

    app.config.globalProperties.$client = window.$client
  },
}

export default plugin
