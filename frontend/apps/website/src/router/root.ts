import { RouteRecordRaw } from 'vue-router'

export const ROOT_PATH = '/'
export const ROOT_PATH_NAME = 'root'

export const rootRoute: RouteRecordRaw = {
  name: ROOT_PATH_NAME,
  path: ROOT_PATH,
  component: () => import('@/layouts/EmptyLayout.vue'),
}
