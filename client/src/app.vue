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
            <div class="settings">
              <neko-menu />
            </div>
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
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    max-width: 100vw;
    height: 100dvh;
    max-height: 100dvh;
    overflow: hidden;
    flex-direction: row;
    display: flex;
    background: #05070c;
    color: $text-normal;

    .neko-main {
      position: relative;
      isolation: isolate;
      width: 100%;
      height: auto;
      min-width: 0;
      max-width: 100%;
      flex: 1 1 auto;
      min-height: 0;
      flex-direction: column;
      display: flex;
      overflow: hidden;
      background: #05070c;

      .header-container {
        position: relative;
        flex: 0 0 $party-header-height;
        z-index: 6;
        width: 100%;
        min-width: 0;
        height: auto;
        display: flex;
        pointer-events: none;
        background: linear-gradient(180deg, rgba(#05070c, 0.86), rgba(#05070c, 0));

        > * {
          pointer-events: auto;
        }
      }

      .video-container {
        position: relative;
        flex: 1 1 auto;
        min-width: 0;
        min-height: 0;
        z-index: 0;
        display: flex;
        padding: $party-gutter;
        background: #05070c;

        .video-surface {
          position: relative;
          flex: 1;
          min-width: 0;
          min-height: 0;
          display: flex;
          overflow: hidden;
          border: 1px solid rgba($text-normal, 0.1);
          border-radius: 20px;
          background: #050a12;
          box-shadow: $elevation-high;

          .video-menu.bottom {
            bottom: $party-gutter;
          }
        }

        &.video-only {
          padding: 0;
          background: #000;

          .video-surface {
            border: 0;
            border-radius: 0;
            box-shadow: none;
          }
        }
      }

      .room-container {
        position: relative;
        flex: 0 0 auto;
        z-index: 5;
        margin: 0 $party-gutter $party-gutter;
        min-height: 0;
        max-width: none;
        display: grid;
        grid-template-rows: auto auto;
        gap: clamp(0.2rem, 0.6vw, 0.45rem);
        padding: clamp(0.45rem, 1vw, 0.75rem) clamp(0.65rem, 1.3vw, 0.875rem);
        overflow: visible;
        border: 1px solid rgba(#fff, 0.12);
        border-radius: 18px;
        background: rgba($background-primary, 0.68);
        box-shadow: 0 18px 50px rgba(0, 0, 0, 0.36);
        backdrop-filter: blur(18px) saturate(130%);

        > .members {
          grid-row: 1;
          height: auto;
          min-height: $party-member-rail-height;
          position: relative;
          z-index: 1;
        }

        .room-menu {
          grid-row: 2;
          max-width: 100%;
          min-height: $party-control-height;
          display: grid;
          grid-template-columns: minmax(8rem, auto) minmax(0, 1fr) minmax(8rem, auto);
          grid-template-areas: 'settings controls emotes';
          align-items: center;
          gap: clamp(0.35rem, 1vw, 0.75rem);
          position: relative;
          z-index: 2;

          .settings {
            grid-area: settings;
            margin-left: 0;
            min-width: 0;
            justify-content: flex-start;
            align-items: center;
            display: flex;
          }

          .controls {
            grid-area: controls;
            min-width: 0;
            min-height: $party-control-height;
            padding: 0 8px;
            border: 1px solid rgba(#fff, 0.08);
            border-radius: 13px;
            background: rgba(#000, 0.22);
            justify-content: center;
            align-items: center;
            display: flex;
            overflow: visible;

            > ul {
              width: 100%;
              display: flex;
              flex: 0 1 auto;
              min-width: 0;
              max-width: 100%;
              flex-wrap: nowrap;
              justify-content: flex-start;
              overflow-x: auto;
              overflow-y: hidden;
              scrollbar-width: thin;

              > li {
                flex: 0 0 auto;
              }
            }
          }

          .emotes {
            grid-area: emotes;
            margin: 0;
            min-width: 0;
            justify-content: flex-end;
            align-items: center;
            display: flex;
          }
        }
      }
    }
  }

  .panel-backdrop {
    position: fixed;
    inset: 0;
    z-index: 9;
    display: block;
    background: rgba(#02040a, 0.58);
    backdrop-filter: blur(4px);
  }

  .app-notifications {
    top: $party-header-height !important;
    pointer-events: none;
  }

  @media only screen and (max-width: 1180px) {
    html,
    body {
      overflow-y: auto !important;
      width: auto !important;
      height: auto !important;
    }

    body > p {
      display: none;
    }

    #neko {
      position: relative;
      flex-direction: column;
      max-height: initial !important;

      .neko-menu {
        position: fixed;
        top: auto;
        right: 0;
        bottom: 0;
        z-index: 12;
        height: min(70dvh, 42rem);
        width: min(100%, 29rem) !important;
        border-radius: 18px 18px 0 0;
      }
    }
  }

  @media only screen and (max-width: 1024px) and (orientation: portrait) {
    #neko {
      &.expanded > .neko-menu {
        height: min(72dvh, 42rem);
        width: min(100%, 29rem) !important;
      }
    }
  }

  @media only screen and (max-width: 768px) {
    #neko {
      .neko-main {
        .header-container {
          flex-basis: $party-header-height;
        }

        .video-container {
          padding: $party-gutter;

          .video-surface {
            border-radius: 14px;

            .video-menu.bottom {
              bottom: $party-gutter;
            }
          }
        }

        .room-container {
          margin: 0 $party-gutter $party-gutter;
          padding: clamp(0.35rem, 1vw, 0.5rem) clamp(0.45rem, 1.2vw, 0.75rem);
          border-radius: 15px;

          > .members {
            flex-basis: auto;
            height: auto;
            min-height: $party-member-rail-height;
          }

          .room-menu {
            min-height: $party-control-height;

            .controls {
              min-height: $party-control-height;
            }

            .emotes {
              min-width: 0;
            }
          }
        }
      }
    }
  }

  @media only screen and (max-width: 480px) {
    #neko .neko-main .room-container .room-menu {
      grid-template-columns: minmax(0, 1fr) minmax(0, 1fr);
      grid-template-areas:
        'settings emotes'
        'controls controls';

      .settings {
        .menu-actions {
          gap: 2px;

          .action-button {
            width: 30px;
            height: 30px;
          }
        }

        .locale-picker select {
          min-width: 48px;
          padding: 0 18px 0 6px;
          font-size: 10px;
        }
      }

      .emotes {
        justify-content: flex-end;
      }
    }
  }

  @media only screen and (max-width: 360px) {
    #neko .neko-main .room-container .room-menu {
      grid-template-columns: 1fr;
      grid-template-areas:
        'settings'
        'controls';

      .emotes {
        display: none;
      }
    }
  }
</style>

<script lang="ts">
  import { Vue, Component, Ref, Watch } from 'vue-property-decorator'

  import Connect from '~/components/connect.vue'
  import Video from '~/components/video.vue'
  import Menu from '~/components/menu.vue'
  import Side from '~/components/side.vue'
  import Controls from '~/components/controls.vue'
  import Members from '~/components/members.vue'
  import Emotes from '~/components/emotes.vue'
  import About from '~/components/about.vue'
  import Header from '~/components/header.vue'
  import Unsupported from '~/components/unsupported.vue'

  @Component({
    name: 'neko',
    components: {
      'neko-connect': Connect,
      'neko-video': Video,
      'neko-menu': Menu,
      'neko-side': Side,
      'neko-controls': Controls,
      'neko-members': Members,
      'neko-emotes': Emotes,
      'neko-about': About,
      'neko-header': Header,
      'neko-unsupported': Unsupported,
    },
  })
  export default class extends Vue {
    @Ref('video') video!: Video

    shakeKbd = false

    get volume() {
      const numberParam = parseFloat(new URL(location.href).searchParams.get('volume') || '1.0')
      return Math.max(0.0, Math.min(!isNaN(numberParam) ? numberParam * 100 : 100, 100))
    }

    get scroll() {
      const numberParam = parseInt(new URL(location.href).searchParams.get('scroll') || '', 10)
      return Math.max(1, Math.min(!isNaN(numberParam) ? numberParam : 10, 100))
    }

    get isCastMode() {
      return !!new URL(location.href).searchParams.get('cast')
    }

    get isEmbedMode() {
      return !!new URL(location.href).searchParams.get('embed')
    }

    get hideControls() {
      return this.isCastMode
    }

    get videoOnly() {
      return this.isCastMode || this.isEmbedMode
    }

    @Watch('volume', { immediate: true })
    onVolume(volume: number) {
      if (new URL(location.href).searchParams.has('volume')) {
        this.$accessor.video.setVolume(volume)
      }
    }

    @Watch('scroll', { immediate: true })
    onScroll(scroll: number) {
      if (new URL(location.href).searchParams.has('scroll')) {
        this.$accessor.settings.setScroll(scroll)
      }
    }

    @Watch('hideControls', { immediate: true })
    onHideControls(enabled: boolean) {
      if (enabled) {
        this.$accessor.video.setMuted(false)
        this.$accessor.settings.setSound(false)
      }
    }

    @Watch('side')
    onSide(side: boolean) {
      if (side) {
        console.log('side enabled')
        // scroll to the side
        this.$nextTick(() => {
          const side = document.querySelector('aside')
          if (side) {
            side.scrollIntoView({ behavior: 'smooth', block: 'start' })
          }
        })
      }
    }

    controlAttempt() {
      if (this.shakeKbd || this.$accessor.remote.hosted) return

      this.shakeKbd = true
      window.setTimeout(() => (this.shakeKbd = false), 5000)
    }

    get about() {
      return this.$accessor.client.about
    }

    get side() {
      return this.$accessor.client.side
    }

    get authenticated() {
      return this.$accessor.connection.authenticated
    }

    closeSide() {
      this.$accessor.client.setSide(false)
    }
  }
</script>
