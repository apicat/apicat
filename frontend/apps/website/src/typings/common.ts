import type { Router } from 'vue-router'
import 'pinia'

declare module 'pinia' {
  export interface PiniaCustomProperties {
    // type the router added by the plugin above (#adding-new-external-properties)
    $router: Router
  }
}

export declare interface Language {
  name: string
  lang: string
}
