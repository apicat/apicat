import { compile } from 'path-to-regexp'
import { PROJECT_DETAIL_PATH } from './constant'
import { RouteRecordRaw } from 'vue-router'
import { MemberAuthorityInProject } from '@/typings/member'

export const SCHEMA_EDIT_NAME = 'definition.schema.edit'
export const SCHEMA_DETAIL_PATH = PROJECT_DETAIL_PATH + '/schema/:schema_id?'

export const SCHEMA_DETAIL_NAME = 'definition.schema.detail'
export const SCHEMA_EDIT_PATH = PROJECT_DETAIL_PATH + '/schema/:schema_id/edit'

export const getSchemaDetailPath = (project_id: number | string, schema_id: number | string) => compile(SCHEMA_DETAIL_PATH)({ project_id, schema_id })
export const getSchemaEditPath = (project_id: number | string, schema_id: number | string) => compile(SCHEMA_EDIT_PATH)({ project_id, schema_id })

export const schemaDetailRoute: RouteRecordRaw = {
  name: SCHEMA_DETAIL_NAME,
  path: SCHEMA_DETAIL_PATH,
  component: () => import('@/views/document/SchemaDetailPage.vue'),
  meta: {
    ignoreAuth: true,
  },
}

export const schemaEditRoute: RouteRecordRaw = {
  name: SCHEMA_EDIT_NAME,
  path: SCHEMA_EDIT_PATH,
  component: () => import('@/views/document/SchemaEditPage.vue'),
  meta: {
    editableRoles: [MemberAuthorityInProject.MANAGER, MemberAuthorityInProject.WRITE],
  },
}
