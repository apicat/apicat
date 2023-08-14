import type { RouteRecordRaw } from 'vue-router'
import { compile } from 'path-to-regexp'
import { ITERATION_DETAIL_PATH, PROJECT_DETAIL_PATH } from './constant'
import { MemberAuthorityInProject } from '@/typings/member'

export const DOCUMENT_DETAIL_NAME = 'document.detail'
export const DOCUMENT_DETAIL_PATH = PROJECT_DETAIL_PATH + '/doc/:doc_id?'

export const ITERATION_ALIAS_DOCUMENT_DETAIL_PATH = ITERATION_DETAIL_PATH + '/doc/:doc_id?'

export const DOCUMENT_EDIT_NAME = 'document.edit'
export const DOCUMENT_EDIT_PATH = PROJECT_DETAIL_PATH + '/doc/:doc_id/edit'

export const getDocumentDetailPath = (project_id: number | string, doc_id: number | string) => compile(DOCUMENT_DETAIL_PATH)({ project_id, doc_id })
export const getDocumentEditPath = (project_id: number | string, doc_id: number | string) => compile(DOCUMENT_DETAIL_PATH)({ project_id, doc_id }) + '/edit'

export const documentDetailRoute: RouteRecordRaw = {
  name: DOCUMENT_DETAIL_NAME,
  path: DOCUMENT_DETAIL_PATH,
  component: () => import('@/views/document/DocumentDetailPage.vue'),
  meta: {
    ignoreAuth: true,
  },
}

export const documentEditRoute: RouteRecordRaw = {
  name: DOCUMENT_EDIT_NAME,
  path: DOCUMENT_EDIT_PATH,
  component: () => import('@/views/document/DocumentEditPage.vue'),
  meta: {
    editableRoles: [MemberAuthorityInProject.MANAGER, MemberAuthorityInProject.WRITE],
  },
}
