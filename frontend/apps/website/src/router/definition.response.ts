import { compile } from 'path-to-regexp'
import { PROJECT_DETAIL_PATH } from './constant'
import { RouteRecordRaw } from 'vue-router'

export const RESPONSE_DETAIL_NAME = 'response.detail'
export const RESPONSE_DETAIL_PATH = PROJECT_DETAIL_PATH + '/response/:response_id?'

export const RESPONSE_EDIT_NAME = 'response.edit'
export const RESPONSE_EDIT_PATH = PROJECT_DETAIL_PATH + '/response/:response_id/edit'

export const getDefinitionResponseDetailPath = (project_id: number | string, response_id: number | string) => compile(RESPONSE_DETAIL_NAME)({ project_id, response_id })
export const getDefinitionResponseEditPath = (project_id: number | string, response_id: number | string) => compile(RESPONSE_EDIT_PATH)({ project_id, response_id }) + '/edit'

const ResponseEditPage = () => import('@/views/definition/response/ResponseEditPage.vue')
const ResponseDetailPage = () => import('@/views/definition/response/ResponseDetailPage.vue')

export const definitionResponseDetailRoute: RouteRecordRaw = {
  name: RESPONSE_DETAIL_NAME,
  path: RESPONSE_DETAIL_PATH,
  component: ResponseDetailPage,
}

export const definitionResponseEditRoute: RouteRecordRaw = {
  name: RESPONSE_EDIT_NAME,
  path: RESPONSE_EDIT_PATH,
  component: ResponseEditPage,
}
