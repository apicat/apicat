import type { RouteRecordRaw } from 'vue-router'
import { USER_PAGE_NAME, USER_PAGE_PATH, USER_SETTING_NAME, USER_SETTING_PATH } from './constant'

const UserSettingPage = () => import('@/views/user/UserSettingPage.vue')

export const userRoute: RouteRecordRaw = {
  path: USER_SETTING_PATH,
  name: USER_SETTING_NAME,
  redirect: {
    name: USER_PAGE_NAME,
    params: {
      page: 'general',
    },
  },
  children: [
    {
      path: USER_PAGE_PATH,
      name: USER_PAGE_NAME,
      component: UserSettingPage,
    },
    {
      path: `${USER_SETTING_PATH}/:path(.*)*`,
      name: 'user.404',
      redirect: {
        name: USER_PAGE_NAME,
        params: {
          page: 'general',
        },
      },
    },
  ],
}
