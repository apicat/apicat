import { INDEX_PATH, LOGIN_PATH, REGISTE_PATH } from './constant'
import AppLayout from '../layout/AppLayout.vue'

const Index = () => import('../views/Index.vue')
const Login = () => import('../views/Login.vue')
const Register = () => import('../views/Register.vue')
const NotFound = () => import('../views/errors/NotFound.vue')

/**
 * index router
 */
export const indexRouter = {
    path: INDEX_PATH,
    name: 'index',
    component: AppLayout,
    meta: { ignoreAuth: true },
    children: [
        {
            path: '',
            name: 'index.feature',
            component: Index,
        },
    ],
}

/**
 * login router
 */
export const loginRouter = {
    path: LOGIN_PATH,
    name: 'login',
    meta: { ignoreAuth: true },
    component: Login,
}

/**
 * register router
 */
export const registerRouter = {
    path: REGISTE_PATH,
    name: 'register',
    meta: { ignoreAuth: true },
    component: Register,
}

/**
 * 404 router
 */
export const notFoundRouter = {
    path: '/:path(.*)*',
    name: 'error',
    component: AppLayout,
    meta: { ignoreAuth: true },
    children: [
        {
            path: '/:path(.*)*',
            name: 'error.not.found',
            component: NotFound,
        },
    ],
}
