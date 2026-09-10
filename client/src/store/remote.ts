import { getterTree, mutationTree, actionTree } from 'typed-vuex'
import { Member } from '~/neko/types'
import { EVENT } from '~/neko/events'
import { accessor } from '~/store'

const keyboardModifierState = (capsLock: boolean, numLock: boolean, scrollLock: boolean) =>
  Number(capsLock) + 2 * Number(numLock) + 4 * Number(scrollLock)

export const namespaced = true

export const state = () => ({
  id: '',
  epoch: 0,
  clipboard: '',
  locked: false,
  // The server sends the effective setting during websocket initialization.
  // Start conservatively so the pre-init state cannot send host-only messages.
  implicitHosting: false,
  fileTransfer: true,
  keyboardModifierState: -1,
})

export const getters = getterTree(state, {
  controlling: (state, getters, root) => {
    return root.user.id === state.id
  },
  hosting: (state, getters, root) => {
    return root.user.id === state.id || state.implicitHosting
  },
  hosted: (state) => {
    return state.id !== '' || state.implicitHosting
  },
  host: (state, getters, root) => {
    return root.user.members[state.id] || (state.implicitHosting && root.user.id) || null
  },
})

export const mutations = mutationTree(state, {
  setHost(state, host: string | Member) {
    if (typeof host === 'string') {
      state.id = host
    } else {
      state.id = host.id
    }
  },

  setEpoch(state, epoch: number) {
    state.epoch = epoch
  },

  setClipboard(state, clipboard: string) {
    state.clipboard = clipboard
  },

  setKeyboardModifierState(state, { capsLock, numLock, scrollLock }) {
    state.keyboardModifierState = keyboardModifierState(capsLock, numLock, scrollLock)
  },

  setLocked(state, locked: boolean) {
    state.locked = locked
  },

  setImplicitHosting(state, val: boolean) {
    state.implicitHosting = val
  },

  setFileTransfer(state, val: boolean) {
    state.fileTransfer = val
  },

  reset(state) {
    state.id = ''
    state.epoch = 0
    state.clipboard = ''
    state.locked = false
    state.implicitHosting = false
    state.keyboardModifierState = -1
  },
})

export const actions = actionTree(
  { state, getters, mutations },
  {
    sendClipboard({ getters }, clipboard: string) {
      // Clipboard writes are host-only on the server. In implicit-hosting
      // mode, `hosting` means that input may request control, not that this
      // session already owns it.
      if (!accessor.connection.connected || !getters.controlling) {
        return
      }

      $client.sendMessage(EVENT.CLIPBOARD.SET, { text: clipboard })
    },

    toggle({ getters }) {
      if (!accessor.connection.connected) {
        return
      }

      if (!getters.hosting) {
        $client.sendMessage(EVENT.CONTROL.REQUEST)
      } else {
        $client.sendMessage(EVENT.CONTROL.RELEASE)
      }
    },

    request({ getters }) {
      if (!accessor.connection.connected || getters.controlling) {
        return
      }

      $client.sendMessage(EVENT.CONTROL.REQUEST)
    },

    release({ getters }) {
      if (!accessor.connection.connected || !getters.hosting) {
        return
      }

      $client.sendMessage(EVENT.CONTROL.RELEASE)
    },

    async give({ getters }, member: string | Member) {
      if (!accessor.connection.connected || !getters.hosting) {
        return
      }

      if (typeof member === 'string') {
        member = accessor.user.members[member]
      }

      if (!member) {
        return
      }

      await $client.room.giveControl(member.id)
    },

    adminControl() {
      if (!accessor.connection.connected || !accessor.user.admin) {
        return
      }

      $client.room.takeControl()
    },

    adminRelease() {
      if (!accessor.connection.connected || !accessor.user.admin) {
        return
      }

      $client.room.resetControl()
    },

    adminGive(store, member: string | Member) {
      if (!accessor.connection.connected) {
        return
      }

      if (typeof member === 'string') {
        member = accessor.user.members[member]
      }

      if (!member) {
        return
      }

      $client.room.giveControl(member.id)
    },

    changeKeyboard({ getters }) {
      if (!accessor.connection.connected || !getters.controlling) {
        return
      }

      $client.sendMessage(EVENT.KEYBOARD.MAP, { layout: accessor.settings.keyboard_layout })
    },

    syncKeyboardModifierState({ state }, { capsLock, numLock, scrollLock }) {
      // The websocket handler rejects this event for viewers. It is common
      // for the pointer to enter the video before control negotiation ends.
      if (!accessor.connection.connected || !accessor.remote.controlling) {
        return
      }

      if (state.keyboardModifierState === keyboardModifierState(capsLock, numLock, scrollLock)) {
        return
      }

      accessor.remote.setKeyboardModifierState({ capsLock, numLock, scrollLock })
      $client.sendMessage(EVENT.KEYBOARD.MODIFIERS, {
        capslock: capsLock,
        numlock: numLock,
      })
    },
  },
)
