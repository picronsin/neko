# Vue 3.5.42 migration

The client now runs on the native Vue 3.5.42 runtime and Vite.

- Vuex and typed-vuex were replaced by Pinia-backed reactive modules.
- Vue 2 plugins were replaced with Vue 3 integrations (`vue-i18n` 11,
  `floating-vue`, `@kyvg/vue3-notification`).
- Components use Vue 3 Options API or the Vue 3-native
  `vue-facing-decorator` package.
- Context and resolution menus are local Vue 3 floating panels.
- Vue CLI, `@vue/compat`, Vuex, and Vue 2 decorator packages are removed.

Build and release gate:

```bash
npm ci
npm run typecheck
npm run test:sdk
npm run build
```

`npm run typecheck` runs `vue-tsc` with the project's current TypeScript and
template settings. Vite does not type-check the production build, so run all
three checks before release.
