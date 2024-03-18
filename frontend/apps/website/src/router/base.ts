import { NO_PERMISSION_PATH } from './constant'

const NotFound = async () => import('@/views/errors/NotFound.vue')
const NoPermission = async () => import('@/views/errors/NoPermission.vue')

export const notFoundRoute = {
  path: '/:path(.*)*',
  name: 'error',
  meta: { ignoreAuth: true, title: 'app.tips.notFound' },
  component: NotFound,
}

export const noPermissionRoute = {
  path: NO_PERMISSION_PATH,
  name: 'no.permission',
  meta: { ignoreAuth: true, title: 'app.tips.noPermission' },
  component: NoPermission,
}
