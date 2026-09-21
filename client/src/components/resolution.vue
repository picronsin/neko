<template>
  <div v-if="visible" class="context" ref="context" :style="menuStyle">
    <li
      v-for="(conf, i) in configurations"
      :key="i"
      @click="screenSet(conf)"
      :class="[conf.width === width && conf.height === height && conf.rate === rate ? 'active' : '']"
    >
      <i class="fas fa-desktop"></i>
      <span>{{ conf.width }}x{{ conf.height }}</span>
      <small>{{ conf.rate }}</small>
    </li>
  </div>
</template>

<style lang="scss" scoped>
  .context {
    background-color: $background-floating;
    background-clip: padding-box;
    border-radius: 0.25rem;
    display: block;
    margin: 0;
    padding: 5px;
    min-width: 150px;
    z-index: 1500;
    position: fixed;
    list-style: none;
    box-sizing: border-box;
    max-height: calc(100% - 50px);
    overflow-y: auto;
    color: $interactive-normal;
    user-select: none;
    box-shadow: $elevation-high;
    scrollbar-width: thin;
    scrollbar-color: $background-secondary transparent;

    &::-webkit-scrollbar {
      width: 8px;
    }

    &::-webkit-scrollbar-track {
      background-color: transparent;
    }

    &::-webkit-scrollbar-thumb {
      background-color: $background-secondary;
      border: 2px solid $background-floating;
      border-radius: 4px;
    }

    &::-webkit-scrollbar-thumb:hover {
      background-color: $background-floating;
    }

    > li {
      margin: 0;
      position: relative;
      align-content: center;
      display: flex;
      flex-direction: row;
      padding: 8px;
      cursor: pointer;
      border-radius: 3px;

      i {
        margin-right: 10px;
      }

      span {
        flex-grow: 1;
      }

      small {
        font-size: 0.7em;
        justify-self: flex-end;
        align-self: flex-end;
      }

      &.active,
      &:hover,
      &:focus {
        text-decoration: none;
        background-color: $background-modifier-hover;
        color: $interactive-hover;
      }

      &:focus {
        outline: 0;
      }
    }

    &:focus {
      outline: 0;
    }
  }
</style>

<script lang="ts">
  import { defineComponent } from 'vue'
  import { ScreenResolution } from '~/neko/types'

  export default defineComponent({
    name: 'neko-resolution',
    components: {},
    data: () => ({ visible: false, menuStyle: {} as Record<string, string> }),
    computed: {
      width() {
        return this.$accessor.video.width
      },

      height() {
        return this.$accessor.video.height
      },

      rate() {
        return this.$accessor.video.rate
      },

      configurations() {
        return this.$accessor.video.configurations
      },
    },
    mounted() {
      // The video overlay deliberately stops pointer events so it can capture
      // remote input. Listen in the capture phase to close the menu even when
      // the click originated on that overlay.
      document.addEventListener('pointerdown', this.onDocumentPointerDown, true)
      document.addEventListener('keydown', this.onDocumentKeyDown)
    },
    beforeUnmount() {
      document.removeEventListener('pointerdown', this.onDocumentPointerDown, true)
      document.removeEventListener('keydown', this.onDocumentKeyDown)
    },
    methods: {
      open(event: MouseEvent) {
        this.visible = true
        this.menuStyle = { left: String(event.clientX) + 'px', top: String(event.clientY) + 'px' }
      },

      onDocumentPointerDown(event: PointerEvent) {
        if (!this.visible) {
          return
        }

        const context = this.$refs.context as HTMLElement | undefined
        if (!context?.contains(event.target as Node)) {
          this.visible = false
        }
      },

      onDocumentKeyDown(event: KeyboardEvent) {
        if (event.key === 'Escape') {
          this.visible = false
        }
      },

      screenSet(resolution: ScreenResolution) {
        this.$accessor.video.screenSet(resolution)
        this.visible = false
      },
    },
  })
</script>
