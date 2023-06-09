import { LOGIN_PATH, NO_PERMISSION_PATH, REGISTER_PATH } from './constant'

export const notFoundRoute = {
  path: '/:path(.*)*',
  name: 'error',
  meta: { ignoreAuth: true },
  component: () => import('@/views/errors/NotFound.vue'),
}

export const noPermissionRoute = {
  path: NO_PERMISSION_PATH,
  name: 'no.permission',
  meta: { ignoreAuth: true },
  component: () => import('@/views/errors/NoPermission.vue'),
}

export const loginRoute = {
  path: LOGIN_PATH,
  name: 'login',
  meta: { ignoreAuth: true },
  component: () => import('@/views/LoginPage.vue'),
}

export const registerRoute = {
  path: REGISTER_PATH,
  name: 'register',
  meta: { ignoreAuth: true },
  component: () => import('@/views/RegisterPage.vue'),
}
