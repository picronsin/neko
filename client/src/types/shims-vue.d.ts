export {}

declare module '*.vue' {
  import type { DefineComponent } from 'vue'
  const component: DefineComponent<Record<string, unknown>, Record<string, unknown>, any>
  export default component
}

declare module '@vue/runtime-core' {
  interface ComponentCustomProperties {
    [key: string]: any
    $accessor: Record<string, any>
    $t: (...args: any[]) => any
    $te: (...args: any[]) => any
    $notify: (...args: any[]) => any
    $http: any
    $client: any
    $log: any
    $swal: (...args: any[]) => any
    $anime: any
  }
}

// Some legacy tooling packages ship an ambient `vue` declaration. Re-export
// the Vue 3 runtime surface so those declarations cannot hide Vue's exports.
declare module 'vue' {
  export * from '@vue/runtime-dom'
}
