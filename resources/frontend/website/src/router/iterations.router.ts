import { RouteRecordRaw } from 'vue-router'

import IterateNav from '@/views/iterations/IterateNav.vue'

const IterateListContainer = () => import('@/views/iterations/IterateListContainer.vue')

export const IterateListRoute = { path: '', name: 'iterations.index', meta: { title: '迭代' }, component: IterateListContainer }

export const IterateRoutes: RouteRecordRaw[] = [IterateListRoute]

export const IterateRootRoute = {
    path: '/iterations',
    name: 'iterations',
    meta: { title: '迭代', activeClassPrefixes: 'iterations', browserTitle: '迭代 - ApiCat' },
    redirect: { name: 'iterations.index' },
}

export default {
    ...IterateRootRoute,
    component: IterateNav,
    children: IterateRoutes,
}
