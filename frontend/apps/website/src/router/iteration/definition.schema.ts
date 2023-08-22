import { compile } from 'path-to-regexp'
import { ITERATION_DETAIL_PATH } from '../constant'
import { RouteRecordRaw } from 'vue-router'
import { MemberAuthorityInProject } from '@/typings/member'
import { resetEmptyPathParams } from '@/commons'

export const ITERATION_SCHEMA_EDIT_NAME = 'iteration.definition.schema.edit'
export const ITERATION_SCHEMA_DETAIL_PATH = ITERATION_DETAIL_PATH + '/schema/:schema_id?'

export const ITERATION_SCHEMA_DETAIL_NAME = 'iteration.definition.schema.detail'
export const ITERATION_SCHEMA_EDIT_PATH = ITERATION_DETAIL_PATH + '/schema/:schema_id/edit'

export const getSchemaDetailPathWithIterationId = (iteration_id: number | string, schema_id: number | string) =>
  compile(ITERATION_SCHEMA_DETAIL_PATH)(resetEmptyPathParams({ iteration_id, schema_id }))

export const getSchemaEditPathWithIterationId = (iteration_id: number | string, schema_id: number | string) =>
  compile(ITERATION_SCHEMA_EDIT_PATH)(resetEmptyPathParams({ iteration_id, schema_id }))

export const iterationSchemaDetailRoute: RouteRecordRaw = {
  name: ITERATION_SCHEMA_DETAIL_NAME,
  path: ITERATION_SCHEMA_DETAIL_PATH,
  component: () => import('@/views/document/SchemaDetailPage.vue'),
  meta: {
    ignoreAuth: true,
  },
}

export const iterationSchemaEditRoute: RouteRecordRaw = {
  name: ITERATION_SCHEMA_EDIT_NAME,
  path: ITERATION_SCHEMA_EDIT_PATH,
  component: () => import('@/views/document/SchemaEditPage.vue'),
  meta: {
    editableRoles: [MemberAuthorityInProject.MANAGER, MemberAuthorityInProject.WRITE],
  },
}
