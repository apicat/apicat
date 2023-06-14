import type { RouteRecordRaw } from 'vue-router'
import { MAIN_PATH_ALIAS, PROJECT_LIST_ROOT_PATH, PROJECT_LIST_ROOT_PATH_NAME } from './constant'

const ProjectListPage = () => import('@/views/project/ProjectListPage.vue')

export const pojectsRoute: RouteRecordRaw = {
  name: PROJECT_LIST_ROOT_PATH_NAME,
  path: PROJECT_LIST_ROOT_PATH,
  alias: MAIN_PATH_ALIAS,
  component: ProjectListPage,
  meta: {},
}
