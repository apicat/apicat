import Node from '@/components/AcTree/model/node'
import 'pinia'
import type { Router } from 'vue-router'
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

/**
 * 当前操作的节点信息
 */
export declare type ActiveNodeInfo = { node: Node | undefined; id: number | undefined }

export declare type ProjectDetailModals = {
  exportDocument: (project_id?: string, doc_id?: string) => void
  shareDocument: (project_id: string, doc_id: string) => void
}
