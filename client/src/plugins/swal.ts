import type { App, Plugin } from 'vue'

import { SweetAlertOptions } from 'sweetalert2'
import Swal from 'sweetalert2/dist/sweetalert2.js'

type VueSwalInstance = typeof Swal.fire

declare module '@vue/runtime-core' {
  interface ComponentCustomProperties {
    $swal: any
  }
}

interface VueSweetalert2Options extends SweetAlertOptions {
  // includeCss?: boolean;
}

const plugin: Plugin = {
  install(app: App, options?: VueSweetalert2Options): void {
    const swalFunction = (...args: [SweetAlertOptions]) => {
      if (options) {
        const mixed = Swal.mixin(options)

        return mixed.fire(...args)
      }

      return Swal.fire(...args)
    }

    let methodName: string | number | symbol

    for (methodName in Swal) {
      // @ts-ignore
      if (Object.prototype.hasOwnProperty.call(Swal, methodName) && typeof Swal[methodName] === 'function') {
        // @ts-ignore
        swalFunction[methodName] = ((method) => {
          return (...args: any[]) => {
            // @ts-ignore
            return Swal[method](...args)
          }
        })(methodName)
      }
    }

    app.config.globalProperties.$swal = swalFunction
  },
}

export default plugin
