import type { RouteRecordRaw } from 'vue-router'
import { compile } from 'path-to-regexp'
import ProjectDetailLayout from '@/layouts/ProjectDetailLayout'

export const PROJECT_DETAIL_PATH = '/proejct/:project_id'

export const DOCUMENT_DETAIL_PATH = PROJECT_DETAIL_PATH + '/doc/:doc_id?'
export const DOCUMENT_EDIT_PATH = PROJECT_DETAIL_PATH + '/doc/:doc_id/edit'

export const SCHEMA_DETAIL_PATH = PROJECT_DETAIL_PATH + '/schema/:shcema_id?'
export const SCHEMA_EDIT_PATH = PROJECT_DETAIL_PATH + '/schema/:shcema_id/edit'

export const DOCUMENT_DETAIL_NAME = 'document.detail'
export const DOCUMENT_EDIT_NAME = 'document.edit'

export const SCHEMA_EDIT_NAME = 'schema.edit'
export const SCHEMA_DETAIL_NAME = 'schema.detail'

export const getProjectDetailPath = (project_id: number | string) => compile(PROJECT_DETAIL_PATH)({ project_id })

export const getDocumentDetailPath = (project_id: number | string, doc_id: number | string) => compile(DOCUMENT_DETAIL_PATH)({ project_id, doc_id })
export const getDocumentEditPath = (project_id: number | string, doc_id: number | string) => compile(DOCUMENT_DETAIL_PATH)({ project_id, doc_id }) + '/edit'

export const getSchemaDetailPath = (project_id: number | string, shcema_id: number | string) => compile(SCHEMA_DETAIL_PATH)({ project_id, shcema_id })
export const getSchemaEditPath = (project_id: number | string, shcema_id: number | string) => compile(SCHEMA_DETAIL_PATH)({ project_id, shcema_id }) + '/edit'

const DocumentDetailPage = () => import('@/views/document/DocumentDetailPage.vue')
const DocumentEditPage = () => import('@/views/document/DocumentEditPage.vue')
const SchemaEditPage = () => import('@/views/document/SchemaEditPage.vue')
const SchemaDetailPage = () => import('@/views/document/SchemaDetailPage.vue')

const documentDetailRoute: RouteRecordRaw = {
  name: DOCUMENT_DETAIL_NAME,
  path: DOCUMENT_DETAIL_PATH,
  component: DocumentDetailPage,
}

const documentEditRoute: RouteRecordRaw = {
  name: DOCUMENT_EDIT_NAME,
  path: DOCUMENT_EDIT_PATH,
  component: DocumentEditPage,
}

const schemaDetailRoute: RouteRecordRaw = {
  name: SCHEMA_DETAIL_NAME,
  path: SCHEMA_DETAIL_PATH,
  component: SchemaDetailPage,
}

const schemaEditRoute: RouteRecordRaw = {
  name: SCHEMA_EDIT_NAME,
  path: SCHEMA_EDIT_PATH,
  component: SchemaEditPage,
}

export const projectDetailRoute: RouteRecordRaw = {
  name: 'project.detail',
  path: PROJECT_DETAIL_PATH,
  component: ProjectDetailLayout,
  children: [documentDetailRoute, documentEditRoute, schemaDetailRoute, schemaEditRoute],
  redirect: { name: DOCUMENT_DETAIL_NAME },
  meta: {},
}
