<template>
  <div class="emotes" @mouseleave="stopSendingEmotes" @mouseup="stopSendingEmotes">
    <ul v-if="!muted">
      <li v-for="emote in recent" :key="emote">
        <div :class="['emote', emote]" @mousedown.stop.prevent="startSendingEmotes(emote)" />
      </li>
      <li class="picker-slot">
        <button
          type="button"
          class="picker-toggle"
          aria-label="表情"
          :aria-expanded="pickerOpen ? 'true' : 'false'"
          @click="togglePicker"
        >
          <i class="fas fa-grin-beam" aria-hidden="true"></i>
        </button>
      </li>
    </ul>
    <div v-show="pickerOpen" class="emote-picker" role="menu">
      <button
        v-for="emote in emotes"
        :key="emote"
        type="button"
        class="emote-option"
        :aria-label="emote"
        @click="selectEmote(emote)"
      >
        <div :class="['emote', emote]" />
      </button>
    </div>
  </div>
</template>

<style lang="scss" scoped>
  .emotes {
    position: relative;

    ul {
      display: flex;
      flex-direction: row;
      justify-content: center;
      align-items: center;

      li {
        font-size: 24px;
        margin: 0 5px;

        i,
        div,
        button {
          cursor: pointer;
        }
      }
    }

    .picker-toggle,
    .emote-option {
      border: 0;
      color: inherit;
      font: inherit;
    }

    .picker-toggle {
      background: transparent;
      padding: 0;
      display: grid;
      place-items: center;
      list-style: none;

      &::-webkit-details-marker {
        display: none;
      }
    }

    .emote-picker {
      background-color: $background-floating;
      background-clip: padding-box;
      border-radius: 0.25rem;
      display: flex;
      margin: 0;
      padding: 5px;
      width: 220px;
      z-index: 1500;
      position: absolute;
      bottom: calc(100% + 0.5rem);
      right: 0;
      box-sizing: border-box;
      max-height: min(340px, 60vh);
      overflow-y: auto;
      color: $interactive-normal;
      flex-wrap: wrap;
      user-select: none;
      box-shadow: $elevation-high;

      .emote-option {
        background: transparent;
        margin: 0;
        position: relative;
        padding: 5px;
        border-radius: 3px;

        .emote {
          width: 24px;
          height: 24px;
        }

        &:hover,
        &:focus-visible {
          text-decoration: none;
          background-color: $background-modifier-hover;
          color: $interactive-hover;
        }

        &:focus-visible {
          outline: 0;
        }
      }

      &:focus {
        outline: 0;
      }
    }
  }
</style>

<script lang="ts">
  import { Vue, Component } from 'vue-property-decorator'
  import { get, set } from '../utils/localstorage'

  @Component({
    name: 'neko-emotes',
  })
  export default class extends Vue {
    recent: string[] = JSON.parse(get('emote_recent', '[]'))
    pickerOpen = false

    get emotes() {
      return [
        'anger',
        'bomb',
        'sleep',
        'explode',
        'sweat',
        'poo',
        'hundred',
        'alert',
        'punch',
        'wave',
        'okay',
        'thumbs-up',
        'clap',
        'prey',
        'celebrate',
        'flame',
        'goof',
        'love',
        'cool',
        'smerk',
        'worry',
        'ouch',
        'cry',
        'surprised',
        'quiet',
        'rage',
        'annoy',
        'steamed',
        'scared',
        'terrified',
        'sleepy',
        'dead',
        'happy',
        'roll-eyes',
        'thinking',
        'clown',
        'sick',
        'rofl',
        'drule',
        'sniff',
        'sus',
        'party',
        'odd',
        'hot',
        'cold',
        'blush',
        'sad',
      ].filter((v) => !this.recent.includes(v))
    }

    get muted() {
      return this.$accessor.user.muted
    }

    sendEmote(emote: string) {
      if (!this.recent.includes(emote)) {
        if (this.recent.length > 4) {
          this.recent.shift()
        }
        this.recent.push(emote)
        set('emote_recent', JSON.stringify(this.recent))
      }
      this.$accessor.chat.sendEmote(emote)
    }

    selectEmote(emote: string) {
      this.sendEmote(emote)
      this.pickerOpen = false
    }

    togglePicker() {
      this.pickerOpen = !this.pickerOpen
    }

    private interval!: number

    startSendingEmotes(emote: string) {
      this.$accessor.chat.sendEmote(emote)
      this.stopSendingEmotes()

      this.interval = window.setInterval(() => {
        this.$accessor.chat.sendEmote(emote)
      }, 350)
    }

    stopSendingEmotes() {
      if (this.interval) {
        clearInterval(this.interval)
      }
    }

  }
</script>
