import type { RouteRecordRaw } from 'vue-router'

export const ROOT_PATH = '/'
export const ROOT_PATH_NAME = 'root'
export const HOME_PATH = '/home'
export const HOME_PATH_NAME = 'home'

export const indexRoute: RouteRecordRaw = {
  name: ROOT_PATH_NAME,
  path: ROOT_PATH,
  redirect: HOME_PATH,
}

export const notFoundRoute = {
  path: '/:path(.*)*',
  name: 'error',
  component: () => import('@/views/errors/NotFound.vue'),
}
