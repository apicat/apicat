import type { RouteRecordRaw } from 'vue-router'
import { compile } from 'path-to-regexp'
import DocumentHistoryLayout from '@/layouts/DocumentHistoryLayout.vue'

const docuemntHistoryPath = '/history/document/:project_id/:doc_id/:history_id?'

export const getDocumentHistoryPath = (project_id: string, doc_id: string, history_id?: string) => compile(docuemntHistoryPath)({ project_id, doc_id, history_id })

// 文档历史记录
export const documentHistoryRoute: RouteRecordRaw = {
  name: 'history.docuemnt',
  path: docuemntHistoryPath,
  component: DocumentHistoryLayout,
  children: [
    {
      name: 'history.docuemnt.detail',
      path: docuemntHistoryPath,
      component: () => import('@/views/document/DocumentHistoryPage.vue'),
    },
  ],
}
