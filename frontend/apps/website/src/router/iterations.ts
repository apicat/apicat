import type { RouteRecordRaw } from 'vue-router'

const IterationsListPage = () => import('@/views/iterations/IterationsListPage.vue')

export const iterationsRoute: RouteRecordRaw = {
  name: 'iterations',
  path: '/iterations',
  component: IterationsListPage,
}
