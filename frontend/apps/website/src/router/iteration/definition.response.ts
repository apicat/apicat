import { compile } from 'path-to-regexp'
import { ITERATION_DETAIL_PATH } from '../constant'
import { RouteRecordRaw } from 'vue-router'
import { MemberAuthorityInProject } from '@/typings/member'

export const ITERATION_RESPONSE_DETAIL_NAME = 'iteration.definition.response.detail'
export const ITERATION_RESPONSE_DETAIL_PATH = ITERATION_DETAIL_PATH + '/response/:response_id?'

export const ITERATION_RESPONSE_EDIT_NAME = 'iteration.definition.response.edit'
export const ITERATION_RESPONSE_EDIT_PATH = ITERATION_DETAIL_PATH + '/response/:response_id/edit'

export const getDefinitionResponseDetailPathWithIterationId = (iteration_id: number | string, response_id: number | string) =>
  compile(ITERATION_RESPONSE_DETAIL_PATH)({ iteration_id, response_id })
export const getDefinitionResponseEditPathWithIterationId = (iteration_id: number | string, response_id: number | string) =>
  compile(ITERATION_RESPONSE_EDIT_PATH)({ iteration_id, response_id })

export const iterationDefinitionResponseDetailRoute: RouteRecordRaw = {
  name: ITERATION_RESPONSE_DETAIL_NAME,
  path: ITERATION_RESPONSE_DETAIL_PATH,
  component: () => import('@/views/definition/response/ResponseDetailPage.vue'),
  meta: {
    ignoreAuth: true,
  },
}

export const iterationDefinitionResponseEditRoute: RouteRecordRaw = {
  name: ITERATION_RESPONSE_EDIT_NAME,
  path: ITERATION_RESPONSE_EDIT_PATH,
  component: () => import('@/views/definition/response/ResponseEditPage.vue'),
  meta: {
    editableRoles: [MemberAuthorityInProject.MANAGER, MemberAuthorityInProject.WRITE],
  },
}
