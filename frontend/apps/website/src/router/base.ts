export const LOGIN_PATH = '/login'
export const REGISTER_PATH = '/register'

export const notFoundRoute = {
  path: '/:path(.*)*',
  name: 'error',
  meta: { ignoreAuth: true },
  component: () => import('@/views/errors/NotFound.vue'),
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
