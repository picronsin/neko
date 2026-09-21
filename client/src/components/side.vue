<template>
  <aside
    class="neko-menu"
    id="room-panel"
    :aria-label="$t('ui.room_panel')"
    @keydown.esc="close"
    @keydown.tab="keepFocus"
  >
    <nav class="tabs-container" :aria-label="$t('ui.room_panel')">
      <button type="button" class="panel-close" :aria-label="$t('ui.close_room_panel')" @click.stop.prevent="close">
        <i class="fas fa-xmark" aria-hidden="true" />
      </button>
      <ul>
        <li>
          <button
            type="button"
            role="tab"
            :aria-selected="tab === 'chat'"
            :class="{ active: tab === 'chat' }"
            @click.stop.prevent="change('chat')"
          >
            <i class="fas fa-comment-alt" aria-hidden="true" />
            <span>{{ $t('side.chat') }}</span>
          </button>
        </li>
        <li v-if="filetransferAllowed">
          <button
            type="button"
            role="tab"
            :aria-selected="tab === 'files'"
            :class="{ active: tab === 'files' }"
            @click.stop.prevent="change('files')"
          >
            <i class="fas fa-file" aria-hidden="true" />
            <span>{{ $t('side.files') }}</span>
          </button>
        </li>
        <li>
          <button
            type="button"
            role="tab"
            :aria-selected="tab === 'settings'"
            :class="{ active: tab === 'settings' }"
            @click.stop.prevent="change('settings')"
          >
            <i class="fas fa-sliders-h" aria-hidden="true" />
            <span>{{ $t('side.settings') }}</span>
          </button>
        </li>
      </ul>
    </nav>
    <div class="page-container">
      <neko-chat v-if="tab === 'chat'" />
      <neko-files v-if="tab === 'files'" />
      <neko-settings v-if="tab === 'settings'" />
    </div>
    <footer class="panel-footer"><neko-menu /></footer>
  </aside>
</template>

<style lang="scss">
  .neko-menu {
    display: flex;
    flex-direction: column;
    min-width: 0;
    min-height: 0;
    overflow: hidden;
    background: $background-primary;
    .tabs-container {
      display: flex;
      align-items: center;
      flex: 0 0 auto;
      padding: 12px;
      gap: 8px;
      border-bottom: 1px solid $background-modifier-accent;
      ul {
        display: flex;
        flex: 1;
        min-width: 0;
        gap: 4px;
      }
      li {
        flex: 1;
      }
      button {
        display: flex;
        align-items: center;
        justify-content: center;
        gap: 7px;
        min-height: 44px;
        width: 100%;
        border: 0;
        border-radius: 8px;
        color: $text-muted;
        background: transparent;
        cursor: pointer;
        font-size: 13px;
      }
      button:hover {
        background: $background-modifier-hover;
        color: $text-normal;
      }
      button.active {
        background: $background-modifier-selected;
        color: $style-primary;
      }
      .panel-close {
        order: 2;
        flex: 0 0 44px;
        width: 44px;
      }
    }
    .page-container {
      display: flex;
      flex: 1;
      min-height: 0;
      overflow: auto;
      padding: 16px;
    }
    .panel-footer {
      padding: 8px 16px max(8px, env(safe-area-inset-bottom));
      border-top: 1px solid $background-modifier-accent;
    }
  }
</style>

<script lang="ts">
  import { defineComponent } from 'vue'

  import Settings from '~/components/settings.vue'
  import Chat from '~/components/chat.vue'
  import Files from '~/components/files.vue'
  import Menu from '~/components/menu.vue'

  export default defineComponent({
    name: 'neko-side',
    components: {
      'neko-settings': Settings,
      'neko-chat': Chat,
      'neko-files': Files,
      'neko-menu': Menu,
    },
    data: () => ({ previousFocus: null as HTMLElement | null }),
    mounted() {
      this.previousFocus = document.activeElement as HTMLElement | null
      this.$nextTick(() => this.$el.querySelector('.panel-close')?.focus())
    },
    beforeUnmount() {
      if (this.previousFocus?.isConnected) this.previousFocus.focus()
    },
    computed: {
      canDownload() {
        return this.$accessor.user.admin || this.$accessor.files.userDownload
      },

      canUpload() {
        return this.$accessor.user.admin || this.$accessor.files.userUpload
      },

      filetransferAllowed() {
        return (
          this.$accessor.remote.fileTransfer &&
          (this.$accessor.user.admin || !this.$accessor.session.isLocked('file_transfer')) &&
          (this.canDownload || this.canUpload)
        )
      },

      tab() {
        return this.$accessor.client.tab
      },
    },
    watch: {
      tab: {
        immediate: true,
        handler() {
          // do not show the files tab if file transfer is disabled
          if (this.tab === 'files' && !this.filetransferAllowed) {
            this.change('chat')
          }
        },
      },

      filetransferAllowed: [
        {
          immediate: true,
          handler() {
            this.onTabChange()
          },
        },
        {
          handler() {
            if (this.filetransferAllowed) this.$accessor.files.refresh()
          },
        },
      ],
    },

    methods: {
      keepFocus(event: KeyboardEvent) {
        if (!window.matchMedia('(max-width: 1100px)').matches) return
        const elements = Array.from(
          (this.$el as HTMLElement).querySelectorAll<HTMLElement>(
            'button:not(:disabled), a[href], input:not(:disabled), select:not(:disabled), textarea:not(:disabled), [tabindex="0"]',
          ),
        ).filter((element) => element.getClientRects().length > 0)
        const first = elements[0]
        const last = elements[elements.length - 1]
        if (event.shiftKey && document.activeElement === first) {
          event.preventDefault()
          last?.focus()
        } else if (!event.shiftKey && document.activeElement === last) {
          event.preventDefault()
          first?.focus()
        }
      },
      onTabChange() {
        if (this.tab === 'files' && !this.filetransferAllowed) this.change('chat')
      },
      change(tab: string) {
        this.$accessor.client.setTab(tab)
      },

      close() {
        this.$accessor.client.setSide(false)
      },
    },
  })
</script>
