<template>
  <div id="neko" class="app-shell" :class="[!videoOnly && side ? 'expanded' : '']">
    <template v-if="!$client.supported">
      <neko-unsupported />
    </template>
    <template v-else>
      <main class="neko-main">
        <div v-if="!videoOnly" class="header-container">
          <neko-header />
        </div>
        <div class="video-container" :class="{ 'video-only': videoOnly }">
          <div class="video-surface">
            <neko-video
              ref="video"
              :hideControls="hideControls"
              :extraControls="isEmbedMode"
              @control-attempt="controlAttempt"
            />
          </div>
        </div>
        <div v-if="!videoOnly" class="room-container">
          <neko-members />
          <div class="room-menu">
            <div class="controls">
              <neko-controls :shakeKbd="shakeKbd" />
            </div>
            <div class="emotes">
              <neko-emotes />
            </div>
          </div>
        </div>
      </main>
      <div v-if="!videoOnly && side" class="panel-backdrop" @click.stop.prevent="closeSide" />
      <neko-side v-if="!videoOnly && side" />
      <neko-connect v-if="!authenticated" />
      <neko-about v-if="about" />
      <notifications
        v-if="!videoOnly"
        group="neko"
        position="top left"
        class="app-notifications"
        :ignoreDuplicates="true"
      />
    </template>
  </div>
</template>

<style lang="scss">
  #neko {
    position: fixed;
    inset: 0;
    display: flex;
    width: 100%;
    height: 100dvh;
    overflow: hidden;
    background: $background-tertiary;
    color: $text-normal;
    .neko-main {
      display: flex;
      flex-direction: column;
      flex: 1;
      min-width: 0;
      min-height: 0;
    }
    .header-container {
      flex: 0 0 calc(64px + env(safe-area-inset-top));
      display: flex;
      z-index: 6;
      padding-top: env(safe-area-inset-top);
    }
    .video-container {
      display: flex;
      flex: 1;
      min-height: 0;
      padding: 0 20px 16px;
      .video-surface {
        position: relative;
        display: flex;
        flex: 1;
        min-width: 0;
        overflow: hidden;
        border-radius: 12px;
        background: #000;
      }
      &.video-only {
        padding: 0;
        .video-surface {
          border-radius: 0;
        }
      }
    }
    .room-container {
      flex: 0 0 auto;
      padding: 0 20px max(12px, env(safe-area-inset-bottom));
      > .members {
        min-height: 48px;
      }
      .room-menu {
        display: flex;
        align-items: center;
        justify-content: space-between;
        gap: 16px;
        border-top: 1px solid $background-modifier-accent;
        padding-top: 8px;
      }
      .controls {
        min-width: 0;
        overflow-x: auto;
        scrollbar-width: thin;
        > ul {
          justify-content: flex-start;
          width: max-content;
        }
      }
      .emotes {
        flex: 0 0 auto;
      }
    }
    > .neko-menu {
      position: relative;
      flex: 0 0 340px;
      width: 340px;
      height: 100%;
      border-left: 1px solid $background-modifier-accent;
    }
    .panel-backdrop {
      display: none;
    }
  }
  .app-notifications {
    top: 64px !important;
    pointer-events: none;
  }
  @media (max-width: 1100px) {
    #neko {
      > .neko-menu {
        position: fixed;
        inset: 0 0 0 auto;
        width: min(380px, 90vw);
        z-index: 12;
      }
      .panel-backdrop {
        display: block;
        position: fixed;
        inset: 0;
        z-index: 11;
        background: rgba(0, 0, 0, 0.5);
      }
    }
  }
  @media (max-width: 600px) {
    #neko {
      .header-container {
        flex-basis: calc(60px + env(safe-area-inset-top));
      }
      .video-container {
        padding: 0 8px 8px;
        .video-surface {
          border-radius: 8px;
        }
      }
      .room-container {
        padding: 0 12px max(8px, env(safe-area-inset-bottom));
        .room-menu {
          gap: 4px;
        }
        .emotes .emotes > ul > li:not(.picker-slot) {
          display: none;
        }
      }
      > .neko-menu {
        inset: auto 0 0;
        width: 100%;
        height: min(78dvh, 680px);
        border: 0;
        border-top: 1px solid $background-modifier-accent;
        border-radius: 18px 18px 0 0;
      }
    }
  }
  @media (max-height: 500px) and (orientation: landscape) {
    #neko {
      .header-container {
        flex-basis: 48px;
      }
      .room-container > .members {
        min-height: 40px;
      }
      .video-container {
        padding-bottom: 4px;
      }
      .room-container {
        padding-bottom: max(4px, env(safe-area-inset-bottom));
      }
    }
  }
</style>

<script lang="ts">
  import { defineComponent } from 'vue'

  import Connect from '~/components/connect.vue'
  import Video from '~/components/video.vue'
  import Side from '~/components/side.vue'
  import Controls from '~/components/controls.vue'
  import Members from '~/components/members.vue'
  import Emotes from '~/components/emotes.vue'
  import About from '~/components/about.vue'
  import Header from '~/components/header.vue'
  import Unsupported from '~/components/unsupported.vue'

  export default defineComponent({
    name: 'neko',
    components: {
      'neko-connect': Connect,
      'neko-video': Video,
      'neko-side': Side,
      'neko-controls': Controls,
      'neko-members': Members,
      'neko-emotes': Emotes,
      'neko-about': About,
      'neko-header': Header,
      'neko-unsupported': Unsupported,
    },
    data: () => ({ shakeKbd: false }),

    computed: {
      volume() {
        const numberParam = parseFloat(new URL(location.href).searchParams.get('volume') || '1.0')
        return Math.max(0.0, Math.min(!isNaN(numberParam) ? numberParam * 100 : 100, 100))
      },

      scroll() {
        const numberParam = parseInt(new URL(location.href).searchParams.get('scroll') || '', 10)
        return Math.max(1, Math.min(!isNaN(numberParam) ? numberParam : 10, 100))
      },

      isCastMode() {
        return !!new URL(location.href).searchParams.get('cast')
      },

      isEmbedMode() {
        return !!new URL(location.href).searchParams.get('embed')
      },

      hideControls() {
        return this.isCastMode
      },

      videoOnly() {
        return this.isCastMode || this.isEmbedMode
      },

      about() {
        return this.$accessor.client.about
      },
      side() {
        return this.$accessor.client.side
      },
      authenticated() {
        return this.$accessor.connection.authenticated
      },
    },

    watch: {
      volume: {
        immediate: true,
        handler(volume: number) {
          if (new URL(location.href).searchParams.has('volume')) {
            this.$accessor.video.setVolume(volume)
          }
        },
      },

      scroll: {
        immediate: true,
        handler(scroll: number) {
          if (new URL(location.href).searchParams.has('scroll')) {
            this.$accessor.settings.setScroll(scroll)
          }
        },
      },

      hideControls: {
        immediate: true,
        handler(enabled: boolean) {
          if (enabled) {
            this.$accessor.video.setMuted(false)
            this.$accessor.settings.setSound(false)
          }
        },
      },
    },

    methods: {
      controlAttempt() {
        if (this.shakeKbd || this.$accessor.remote.hosted) return

        this.shakeKbd = true
        window.setTimeout(() => (this.shakeKbd = false), 5000)
      },

      closeSide() {
        this.$accessor.client.setSide(false)
      },
    },
  })
</script>
