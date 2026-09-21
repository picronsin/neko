<template>
  <div class="clipboard" v-if="opened" @click="$event.stopPropagation()">
    <textarea ref="textarea" v-model="clipboard" @focus="selectText" />
  </div>
</template>

<style lang="scss" scoped>
  .clipboard {
    background-color: $background-primary;
    border-radius: 0.25rem;
    display: block;
    padding: 5px;

    position: absolute;
    bottom: clamp(0.5rem, 1.5vw, 1rem);
    right: clamp(0.5rem, 1.5vw, 1rem);

    &,
    textarea {
      width: min(20rem, calc(100vw - 1rem));
      max-width: calc(100vw - 1rem);
      height: min(8rem, 35dvh);
      max-height: calc(100dvh - 1rem);
    }

    textarea {
      border: 0;
      color: $text-normal;
      background: none;

      &::selection {
        background: $text-normal;
      }
    }
  }
</style>

<script lang="ts">
  import { defineComponent } from 'vue'

  export default defineComponent({
    name: 'neko-clipboard',
    data: () => ({ opened: false, typing: undefined as number | undefined }),
    computed: {
      clipboard: {
        get() {
          return this.$accessor.remote.clipboard
        },
        set(data: string) {
          this.$accessor.remote.setClipboard(data)
          if (this.typing) clearTimeout(this.typing)
          this.typing = window.setTimeout(() => this.$accessor.remote.sendClipboard(this.clipboard), 500)
        },
      },
    },
    methods: {
      open() {
        this.opened = true
        document.body.addEventListener('click', this.close)
        window.setTimeout(() => (this.$refs.textarea as HTMLTextAreaElement).focus(), 0)
      },
      selectText(event: FocusEvent) {
        ;(event.target as HTMLTextAreaElement | null)?.select()
      },
      close() {
        this.opened = false
        document.body.removeEventListener('click', this.close)
      },
    },
    beforeUnmount() {
      this.close()
      if (this.typing) clearTimeout(this.typing)
    },
  })
</script>
