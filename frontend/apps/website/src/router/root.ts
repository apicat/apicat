import { RouteRecordRaw } from 'vue-router'
import { ROOT_PATH, ROOT_PATH_NAME } from './constant'

export const rootRoute: RouteRecordRaw = {
  name: ROOT_PATH_NAME,
  path: ROOT_PATH,
  component: () => import('@/layouts/EmptyLayout.vue'),
}
