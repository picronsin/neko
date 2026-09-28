<template>
  <section class="voice-call" :aria-label="$t('voice.title')">
    <button v-if="!joined" data-testid="voice-join" :disabled="busy || !connected" @click="join">
      {{ $t(busy ? 'voice.connecting' : 'voice.join') }}
    </button>
    <template v-else>
      <button data-testid="voice-mute" :aria-pressed="muted" @click="toggleMute">
        {{ $t(muted ? 'ui.unmute' : 'ui.mute') }}
      </button>
      <button data-testid="voice-leave" @click="leave">{{ $t('voice.leave') }}</button>
      <ul>
        <li v-for="member in members" :key="member.id">{{ member.name }}{{ member.muted ? ' 🔇' : '' }}</li>
      </ul>
    </template>
    <p v-if="error" role="alert">{{ error }}</p>
  </section>
</template>

<script lang="ts">
  import { defineComponent, markRaw } from 'vue'
  import { VoiceCall, VoiceMember } from '~/sdk/voice'

  export default defineComponent({
    name: 'neko-voice',
    data() {
      return {
        call: null as VoiceCall | null,
        members: [] as VoiceMember[],
        busy: false,
        joined: false,
        muted: false,
        error: '',
      }
    },
    computed: {
      connected(): boolean {
        return this.$accessor.connection.connected
      },
    },
    watch: {
      connected(value: boolean) {
        if (!value) void this.leave()
      },
    },
    beforeUnmount() {
      void this.call?.leave()
    },
    methods: {
      async join() {
        this.busy = true
        this.error = ''
        if (!navigator.mediaDevices?.getUserMedia) {
          this.error = this.$t('voice.secure_required')
          this.busy = false
          return
        }
        const call = markRaw(
          new VoiceCall(
            this.$http,
            (members) => {
              this.members = members
            },
            (error) => {
              this.error = error.message
              this.joined = false
            },
          ),
        )
        this.call = call
        try {
          await call.join()
          this.joined = this.connected
          this.muted = false
        } catch (error) {
          this.error = error instanceof Error ? error.message : String(error)
        } finally {
          this.busy = false
        }
      },
      async leave() {
        this.joined = false
        await this.call?.leave()
      },
      toggleMute() {
        this.muted = !this.muted
        this.call?.setMuted(this.muted)
      },
    },
  })
</script>

<style scoped>
  .voice-call {
    padding: 10px;
    border-bottom: 1px solid #444;
    flex-shrink: 0;
  }
  button {
    padding: 8px;
    margin-right: 6px;
    border-radius: 6px;
    cursor: pointer;
  }
  p {
    color: #f88;
    font-size: 12px;
  }
  li {
    font-size: 12px;
    padding: 3px;
  }
</style>
