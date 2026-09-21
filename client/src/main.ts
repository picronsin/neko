import './assets/styles/main.scss'

import { createApp } from 'vue'
import Notifications from '@kyvg/vue3-notification'
import FloatingVue from 'floating-vue'
import Logger from './plugins/log'
import Client from './plugins/neko'
import Axios from './plugins/axios'
import Swal from './plugins/swal'
import Anime from './plugins/anime'
import { createVueNekoRuntime } from './neko/vue-adapter'

import { i18n } from './plugins/i18n'
import store, { accessor } from './store'
import app from './app.vue'

const application = createApp(app)
application.use(store)
application.use(i18n)
application.use(Logger)
application.use(Notifications)
application.use(FloatingVue)
application.use(Axios)
application.use(Swal)
application.use(Anime)
application.use(Client)
application.config.globalProperties.$accessor = accessor

window.$client.init(
  createVueNekoRuntime(
    application.config.globalProperties as any,
    (key, params) => i18n.global.t(key, params as any) as string,
  ),
)
accessor.initialise()
application.mount('#neko')
