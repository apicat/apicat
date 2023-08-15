import type { RouteRecordRaw } from 'vue-router'
import { compile } from 'path-to-regexp'
import { ITERATION_DETAIL_PATH } from '../constant'
import { MemberAuthorityInProject } from '@/typings/member'

export const ITERATION_DOCUMENT_DETAIL_NAME = 'iteration.document.detail'
export const ITERATION_DOCUMENT_DETAIL_PATH = ITERATION_DETAIL_PATH + '/doc/:doc_id?'

export const ITERATION_DOCUMENT_EDIT_NAME = 'iteration.document.edit'
export const ITERATION_DOCUMENT_EDIT_PATH = ITERATION_DETAIL_PATH + '/doc/:doc_id/edit'

export const getDocumentDetailPathWithIterationId = (iteration_id: number | string, doc_id: number | string) => compile(ITERATION_DOCUMENT_DETAIL_PATH)({ iteration_id, doc_id })
export const getDocumentEditPathWithIterationId = (iteration_id: number | string, doc_id: number | string) =>
  compile(ITERATION_DOCUMENT_DETAIL_PATH)({ iteration_id, doc_id }) + '/edit'

export const iterationDocumentDetailRoute: RouteRecordRaw = {
  name: ITERATION_DOCUMENT_DETAIL_NAME,
  path: ITERATION_DOCUMENT_DETAIL_PATH,
  component: () => import('@/views/document/DocumentDetailPage.vue'),
  meta: {
    ignoreAuth: true,
  },
}

export const iterationDocumentEditRoute: RouteRecordRaw = {
  name: ITERATION_DOCUMENT_EDIT_NAME,
  path: ITERATION_DOCUMENT_EDIT_PATH,
  component: () => import('@/views/document/DocumentEditPage.vue'),
  meta: {
    editableRoles: [MemberAuthorityInProject.MANAGER, MemberAuthorityInProject.WRITE],
  },
}
