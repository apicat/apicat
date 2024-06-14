import type { RouteRecordRaw } from 'vue-router'
import { compile } from 'path-to-regexp'
import CollectionHistoryLayout from '@/layouts/CollectionHistoryLayout.vue'
import SchemaHistoryLayout from '@/layouts/SchemaHistoryLayout.vue'

const collectionHistoryPath = '/projects/:projectID/collection/:collectionID/history/:historyID?'
const shcemaHistoryPath = '/projects/:projectID/model/:schemaID/history/:historyID?'

export function getCollectionHistoryPath(projectID: string, collectionID: number, historyID?: number) {
  return compile(collectionHistoryPath)({ projectID, collectionID, historyID })
}

export function getSchemaHistoryPath(projectID: string, schemaID: number, historyID?: number) {
  return compile(shcemaHistoryPath)({ projectID, schemaID, historyID })
}

// 文档历史记录
export const collectionHistoryRoute: RouteRecordRaw = {
  name: 'history.collection',
  path: collectionHistoryPath,
  component: CollectionHistoryLayout,
  redirect: { name: 'history.collection.detail' },
  props: true,
  children: [
    {
      name: 'history.collection.detail',
      path: collectionHistoryPath,
      component: () => import('@/views/collection/CollectionHistoryPage.vue'),
      props: true,
      meta: { title: 'app.pageTitles.collectionHistory' },
    },
  ],
}

// 模型历史记录
export const schemaHistoryRoute: RouteRecordRaw = {
  name: 'history.schema',
  path: shcemaHistoryPath,
  component: SchemaHistoryLayout,
  redirect: { name: 'history.schema.detail' },
  props: true,
  children: [
    {
      name: 'history.schema.detail',
      path: shcemaHistoryPath,
      component: () => import('@/views/schema/SchemaHistoryPage.vue'),
      props: true,
      meta: { title: 'app.pageTitles.schemaHistory' },
    },
  ],
}
