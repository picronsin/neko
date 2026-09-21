import type { Plugin } from 'vue'
import axios, { AxiosStatic } from 'axios'

declare global {
  const $http: AxiosStatic

  interface Window {
    $http: any
  }
}

declare module '@vue/runtime-core' {
  interface ComponentCustomProperties {
    $http: AxiosStatic
  }
}

const plugin: Plugin = {
  install(app) {
    window.$http = axios
    app.config.globalProperties.$http = window.$http
  },
}

export default plugin
