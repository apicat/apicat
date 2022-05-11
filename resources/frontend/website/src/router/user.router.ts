import UserNav from '@/views/user/UserNav.vue'
import { RouteRecordRaw } from 'vue-router'
const UserProfile = () => import('../views/user/UserProfile.vue')
const ModifyPassword = () => import('../views/user/ModifyPassword.vue')

export const UserRoutes: RouteRecordRaw[] = [
    { path: 'profile', name: 'user.profile', meta: { title: '个人信息' }, component: UserProfile },
    { path: 'password', name: 'user.password', meta: { title: '修改密码' }, component: ModifyPassword },
]

export default {
    path: '/user',
    name: 'user',
    redirect: { name: 'user.profile' },
    component: UserNav,
    children: UserRoutes,
}
