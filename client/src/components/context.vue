<template>
  <div v-if="contextData" class="context" ref="context" :style="menuStyle" @click.stop>
    <li class="header">
      <div class="user">
        <neko-avatar class="avatar" :seed="contextData.member.displayname" :size="25" />
        <strong>{{ contextData.member.displayname }}</strong>
      </div>
    </li>
    <li class="seperator" />
    <li>
      <span @click="ignore(contextData.member)" v-if="!contextData.member.ignored">{{ $t('context.ignore') }}</span>
      <span @click="unignore(contextData.member)" v-else>{{ $t('context.unignore') }}</span>
    </li>

    <template v-if="admin">
      <li>
        <span @click="mute(contextData.member)" v-if="!contextData.member.muted">{{ $t('context.mute') }}</span>
        <span @click="unmute(contextData.member)" v-else>{{ $t('context.unmute') }}</span>
      </li>
      <li v-if="contextData.member.id === host && !implicitHosting">
          <span @click="adminRelease()">{{ $t('context.release') }}</span>
      </li>
      <li v-if="contextData.member.id === host && !implicitHosting">
          <span @click="adminControl()">{{ $t('context.take') }}</span>
      </li>
      <li>
        <span v-if="contextData.member.id !== host && !implicitHosting" @click="adminGive(contextData.member)">{{
          $t('context.give')
        }}</span>
      </li>
    </template>
    <template v-else>
      <li v-if="hosting && !implicitHosting">
        <span @click="give(contextData.member)">{{ $t('context.give') }}</span>
      </li>
    </template>

    <template v-if="admin && !contextData.member.admin">
      <li class="seperator" />
      <li>
        <span class="danger-action" @click="kick(contextData.member)">{{ $t('context.kick') }}</span>
      </li>
      <li>
        <span class="danger-action" @click="ban(contextData.member)">{{ $t('context.ban') }}</span>
      </li>
    </template>
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

    > li {
      margin: 0;
      position: relative;
      align-content: center;

      &.header {
        .user {
          display: flex;
          flex-direction: row;
          align-content: center;
          padding: 5px 0;

          .avatar {
            width: 25px;
            height: 25px;
            border-radius: 50%;
            margin-right: 5px;
          }

          strong {
            line-height: 25px;
            font-weight: 700;
            max-width: 200px;
            text-overflow: ellipsis;
          }
        }
      }

      &.seperator {
        height: 1px;
        background: $background-secondary;
        margin: 3px 0;
      }

      > span {
        cursor: pointer;
        display: block;
        padding: 5px;
        font-weight: 400;
        text-decoration: none;
        white-space: nowrap;
        background-color: transparent;
        border-radius: 3px;

        &:hover,
        &:focus {
          text-decoration: none;
          background-color: $background-modifier-hover;
          color: $interactive-hover;
        }

        &:focus {
          outline: 0;
        }

        &.danger-action {
          color: $style-error;
        }
      }
    }

    &:focus {
      outline: 0;
    }
  }
</style>

<script lang="ts">
  import { defineComponent } from 'vue'
  import { Member } from '~/neko/types'

  import Avatar from './avatar.vue'

  export default defineComponent({
    name: 'neko-context',
    components: {
      'neko-avatar': Avatar,
    },
    data: () => ({ contextData: null as any, menuStyle: {} as Record<string, string> }),
    computed: {
      admin() {
        return this.$accessor.user.admin
      },

      hosting() {
        return this.$accessor.remote.hosting
      },

      host() {
        return this.$accessor.remote.id
      },

      implicitHosting() {
        return this.$accessor.remote.implicitHosting
      },
    },
    methods: {
      open(event: MouseEvent, data: any) {
        this.contextData = data
        this.menuStyle = { left: String(event.clientX) + 'px', top: String(event.clientY) + 'px' }
      },

      async kick(member: Member) {
        const value = await this.$swal({
          title: this.$t('context.confirm.kick_title', { name: member.displayname }) as string,
          text: this.$t('context.confirm.kick_text', { name: member.displayname }) as string,
          icon: 'warning',
          showCancelButton: true,
          confirmButtonText: this.$t('context.confirm.button_yes') as string,
          cancelButtonText: this.$t('context.confirm.button_cancel') as string,
        })

        if (value) {
          this.$accessor.user.kick(member)
        }
      },

      async ban(member: Member) {
        const value = await this.$swal({
          title: this.$t('context.confirm.ban_title', { name: member.displayname }) as string,
          text: this.$t('context.confirm.ban_text', { name: member.displayname }) as string,
          icon: 'warning',
          showCancelButton: true,
          confirmButtonText: this.$t('context.confirm.button_yes') as string,
          cancelButtonText: this.$t('context.confirm.button_cancel') as string,
        })

        if (value) {
          this.$accessor.user.ban(member)
        }
      },

      async mute(member: Member) {
        const value = await this.$swal({
          title: this.$t('context.confirm.mute_title', { name: member.displayname }) as string,
          text: this.$t('context.confirm.mute_text', { name: member.displayname }) as string,
          icon: 'warning',
          showCancelButton: true,
          confirmButtonText: this.$t('context.confirm.button_yes') as string,
          cancelButtonText: this.$t('context.confirm.button_cancel') as string,
        })

        if (value) {
          this.$accessor.user.mute(member)
        }
      },

      async unmute(member: Member) {
        const value = await this.$swal({
          title: this.$t('context.confirm.unmute_title', { name: member.displayname }) as string,
          text: this.$t('context.confirm.unmute_text', { name: member.displayname }) as string,
          icon: 'warning',
          showCancelButton: true,
          confirmButtonText: this.$t('context.confirm.button_yes') as string,
          cancelButtonText: this.$t('context.confirm.button_cancel') as string,
        })

        if (value) {
          this.$accessor.user.unmute(member)
        }
      },

      adminRelease() {
        this.$accessor.remote.adminRelease()
      },

      adminControl() {
        this.$accessor.remote.adminControl()
      },

      adminGive(member: Member) {
        this.$accessor.remote.adminGive(member)
      },

      give(member: Member) {
        this.$accessor.remote.give(member)
      },

      ignore(member: Member) {
        this.$accessor.user.setIgnored({ id: member.id, ignored: true })
      },

      unignore(member: Member) {
        this.$accessor.user.setIgnored({ id: member.id, ignored: false })
      },
    },
  })
</script>
