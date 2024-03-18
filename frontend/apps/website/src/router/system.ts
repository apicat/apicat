import type { RouteRecordRaw } from 'vue-router'
import { SYSTEM_PAGE_NAME, SYSTEM_PAGE_PATH, SYSTEM_SETTING_NAME, SYSTEM_SETTING_PATH } from './constant'

const SystemSettingPage = () => import('@/views/system/SystemSettingPage.vue')

export const systemRoute: RouteRecordRaw = {
  path: SYSTEM_SETTING_PATH,
  name: SYSTEM_SETTING_NAME,
  redirect: {
    name: SYSTEM_PAGE_NAME,
    params: {
      page: 'general',
    },
  },
  children: [
    {
      path: SYSTEM_PAGE_PATH,
      name: SYSTEM_PAGE_NAME,
      component: SystemSettingPage,
    },
    {
      path: `${SYSTEM_SETTING_PATH}/:path(.*)*`,
      name: 'user.404',
      redirect: {
        name: SYSTEM_PAGE_NAME,
        params: {
          page: 'general',
        },
      },
    },
  ],
}
