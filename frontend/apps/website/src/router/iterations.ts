import type { RouteRecordRaw } from 'vue-router'
import { ITERATION_LIST_ROOT_PATH, ITERATION_LIST_ROOT_PATH_NAME } from './constant'

const IterationsListPage = () => import('@/views/iterations/IterationsListPage.vue')

export const iterationsRoute: RouteRecordRaw = {
  name: ITERATION_LIST_ROOT_PATH_NAME,
  path: ITERATION_LIST_ROOT_PATH,
  component: IterationsListPage,
  meta: {
    title: '迭代列表',
  },
}
