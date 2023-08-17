import type { RouteRecordRaw } from 'vue-router'
import { compile } from 'path-to-regexp'
import DocumentHistoryLayout from '@/layouts/DocumentHistoryLayout.vue'
import SchemaHistoryLayout from '@/layouts/SchemaHistoryLayout.vue'
import { resetEmptyPathParams } from '@/commons'

const docuemntHistoryPath = '/history/document/:project_id/:doc_id/:history_id?'
const shcemaHistoryPath = '/history/schema/:project_id/:schema_id/:history_id?'

export const getDocumentHistoryPath = (project_id: string, doc_id: string, history_id?: string) =>
  compile(docuemntHistoryPath)(resetEmptyPathParams({ project_id, doc_id, history_id }))

export const getSchemaHistoryPath = (project_id: string, schema_id: string, history_id?: string) =>
  compile(shcemaHistoryPath)(resetEmptyPathParams({ project_id, schema_id, history_id }))

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

// 模型历史记录
export const schemaHistoryRoute: RouteRecordRaw = {
  name: 'history.schema',
  path: shcemaHistoryPath,
  component: SchemaHistoryLayout,
  children: [
    {
      name: 'history.schema.detail',
      path: shcemaHistoryPath,
      component: () => import('@/views/definition/schema/SchemaHistoryPage.vue'),
    },
  ],
}
