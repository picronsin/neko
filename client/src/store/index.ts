import { createPinia, defineStore, setActivePinia } from 'pinia'
import { reactive } from 'vue'
import { mutationTree, getterTree, actionTree } from './helpers'
import { readonlyState } from './readonly-state'

import * as video from './video'
import * as chat from './chat'
import * as files from './files'
import * as openinapp from './openinapp'
import * as remote from './remote'
import * as user from './user'
import * as settings from './settings'
import * as client from './client'
import * as emoji from './emoji'
import * as connection from './connection'
import * as session from './session'

export const state = () => ({})

export const mutations = mutationTree(state, {})

export const getters = getterTree(state, {})

export const actions = actionTree(
  { state, getters, mutations },
  {
    initialise() {
      accessor.emoji.initialise()
      accessor.settings.initialise()
    },
  },
)

export const storePattern = {
  state,
  mutations,
  actions,
  getters,
  modules: { connection, video, chat, files, openinapp, user, remote, settings, client, emoji, session },
}

type StoreModule = {
  state: () => Record<string, unknown>
  mutations: Record<string, (state: any, ...args: any[]) => unknown>
  getters: Record<string, (...args: any[]) => unknown>
  actions: Record<string, (context: any, ...args: any[]) => unknown>
}

const modules = storePattern.modules as unknown as Record<string, StoreModule>

type Commands<T> = {
  readonly [K in keyof T]: T[K] extends (context: any, ...args: infer A) => infer R ? (...args: A) => R : never
}

type ModuleAccessor<T extends StoreModule> = Readonly<ReturnType<T['state']>> &
  Commands<T['mutations']> &
  Commands<T['actions']> & {
    readonly [K in keyof T['getters']]: ReturnType<T['getters'][K]>
  }

export type NekoAccessor = {
  readonly [K in keyof typeof storePattern.modules]: ModuleAccessor<(typeof storePattern.modules)[K]>
} & { initialise(): void }

// Pinia owns the reactive module state.  `accessor` deliberately preserves the
// public API used throughout the client during the component migration; it is
// not a Vuex compatibility layer and has no Vuex dependency.
export const pinia = createPinia()
setActivePinia(pinia)

const useNekoStore = defineStore('neko', () => {
  const state = reactive<Record<string, Record<string, unknown>>>({})
  for (const [name, definition] of Object.entries(modules)) {
    state[name] = reactive(definition.state())
  }
  return { state }
})

const piniaStore = useNekoStore()
const moduleAccessors: Record<string, Record<string, unknown>> = {}

for (const [name, definition] of Object.entries(modules)) {
  const moduleState = piniaStore.state[name]
  const getters = new Proxy({} as Record<string, unknown>, {
    get(_target, property: string) {
      const getter = definition.getters[property]
      // Getter implementations use the third argument as the root module map
      // (for example, remote getters read `root.user`).  `piniaStore.state`
      // already is that map; accessing `.state` again produced `undefined`.
      return getter ? (getter as any)(moduleState, getters, piniaStore.state) : undefined
    },
  })

  moduleAccessors[name] = new Proxy({} as Record<string, unknown>, {
    get(_target, property: string) {
      if (property in definition.mutations) {
        return (...args: unknown[]) => definition.mutations[property](moduleState, ...args)
      }
      if (property in definition.actions) {
        return (...args: unknown[]) => definition.actions[property]({ state: moduleState, getters }, ...args)
      }
      if (property in definition.getters) {
        return readonlyState(getters[property])
      }
      return readonlyState(moduleState[property])
    },
    set(_target, property: string) {
      throw new TypeError(`Cannot write ${name}.${property} directly; use a store command`)
    },
    deleteProperty() {
      throw new TypeError('Use a store command to remove state')
    },
    defineProperty() {
      throw new TypeError('Use a store command to define state')
    },
  })
}

export const accessor = {
  ...moduleAccessors,
  initialise() {
    ;(moduleAccessors.emoji.initialise as () => void)()
    ;(moduleAccessors.settings.initialise as () => void)()
  },
} as NekoAccessor

export default pinia
