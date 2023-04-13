import { pojectsRoute } from './projects'
import { projectDetailRoute } from './document'

import { RouteRecordRaw, createRouter, createWebHistory } from 'vue-router'

export * from './document'

export const ROOT_PATH = '/'
export const ROOT_PATH_NAME = 'root'
export const HOME_PATH = '/home'
export const HOME_PATH_NAME = 'home'

const indexRoute: RouteRecordRaw = {
  name: ROOT_PATH_NAME,
  path: ROOT_PATH,
  redirect: HOME_PATH,
}

const notFoundRoute = {
  path: '/:path(.*)*',
  name: 'error',
  component: () => import('@/views/errors/NotFound.vue'),
}

export const router = createRouter({
  history: createWebHistory(),
  routes: [indexRoute, pojectsRoute, projectDetailRoute, notFoundRoute],
})

// 路由拦截器
export * from './filter'

export default router
