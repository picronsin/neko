<template>
  <ul class="video-menu top">
    <li>
      <button
        type="button"
        class="video-action"
        :aria-label="$t('ui.enter_fullscreen')"
        @click.stop.prevent="$emit('requestFullscreen')"
      >
        <i class="fas fa-expand" aria-hidden="true" />
      </button>
    </li>
    <li v-if="admin">
      <button
        type="button"
        class="video-action"
        :aria-label="$t('ui.change_resolution')"
        @click.stop.prevent="$emit('openResolution', $event)"
      >
        <i class="fas fa-desktop" aria-hidden="true" />
      </button>
    </li>
    <li v-if="!controlLocked && !implicitHosting" :class="extraControls ? '' : 'extra-control'">
      <button
        type="button"
        class="video-action"
        :class="[hosted && !hosting ? 'disabled' : '', !hosted && !hosting ? 'faded' : '']"
        :aria-label="$t('ui.request_or_release_control')"
        @click.stop.prevent="$emit('toggleControl')"
      >
        <i class="fas fa-computer-mouse" aria-hidden="true" />
      </button>
    </li>
  </ul>
  <ul class="video-menu bottom">
    <li v-if="hosting && (!clipboardReadAvailable || !clipboardWriteAvailable)">
      <button
        type="button"
        class="video-action"
        :aria-label="$t('ui.open_clipboard')"
        @click.stop.prevent="$emit('openClipboard')"
      >
        <i class="fas fa-clipboard" aria-hidden="true" />
      </button>
    </li>
    <li>
      <button
        type="button"
        class="video-action"
        v-if="pipAvailable"
        @click.stop.prevent="$emit('requestPictureInPicture')"
        v-tooltip="{
          content: $t('ui.picture_in_picture'),
          placement: 'left',
          offset: 5,
          boundariesElement: 'body',
        }"
        :aria-label="$t('ui.picture_in_picture')"
      >
        <i class="fas fa-external-link-alt" aria-hidden="true" />
      </button>
    </li>
    <li v-if="hosting && touchDevice" :class="extraControls ? '' : 'extra-control'">
      <button
        type="button"
        class="video-action"
        :aria-label="$t('ui.open_keyboard')"
        @click.stop.prevent="$emit('openMobileKeyboard')"
      >
        <i class="fas fa-keyboard" aria-hidden="true" />
      </button>
    </li>
  </ul>
</template>

<script setup lang="ts">
  defineProps<{
    admin: boolean
    controlLocked: boolean
    implicitHosting: boolean
    extraControls: boolean
    hosted: boolean
    hosting: boolean
    clipboardReadAvailable: boolean
    clipboardWriteAvailable: boolean
    pipAvailable: boolean
    touchDevice: boolean
  }>()

  defineEmits([
    'requestFullscreen',
    'openResolution',
    'toggleControl',
    'openClipboard',
    'requestPictureInPicture',
    'openMobileKeyboard',
  ])
</script>

<style lang="scss" scoped>
  .video-menu {
    position: absolute;
    z-index: 7;
    right: $party-gutter;

    &.top {
      top: $party-gutter;
    }

    &.bottom {
      bottom: 15px;
    }

    li {
      margin: 0 0 10px 0;

      .video-action {
        display: grid;
        place-items: center;
        width: 30px;
        height: 30px;
        padding: 0;
        border: 1px solid rgba(#fff, 0.12);
        background: rgba(#fff, 0.12);
        border-radius: 9px;
        font-size: 14px;
        text-align: center;
        color: rgba($color: #fff, $alpha: 0.6);
        cursor: pointer;
        transition: color 0.18s ease, background 0.18s ease, border-color 0.18s ease, transform 0.18s ease;

        &:hover,
        &:focus-visible {
          color: #fff;
          background: rgba(#fff, 0.22);
          border-color: rgba(#fff, 0.28);
          transform: translateY(-1px);
        }

        &.faded,
        &.faded i {
          color: rgba($color: $text-normal, $alpha: 0.4);
        }

        &.disabled,
        &.disabled i {
          color: rgba($color: $style-error, $alpha: 0.4);
        }
      }

      &:last-child {
        margin: 0;
      }
    }
  }
</style>
