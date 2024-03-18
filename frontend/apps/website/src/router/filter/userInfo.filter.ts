import { isEmpty } from 'lodash-es'
import type { Router } from 'vue-router'
import { storeToRefs } from 'pinia'
import { systemRoute } from '../system'
import { MAIN_PATH_NAME, SYSTEM_SETTING_NAME } from '../constant'
import { useUserStore } from '@/store/user'

export function setupGetUserInfoFilter(router: Router) {
  router.beforeEach(async (to, from, next) => {
    const userStore = useUserStore()
    const { userInfo, isLogin, isAdmin } = storeToRefs(userStore)
    // 检查用户信息
    if (isLogin.value && (!userInfo.value || isEmpty(userInfo.value)))
      await userStore.getUserInfo()

    // 添加系统路由
    if (isLogin.value && isAdmin.value && !router.hasRoute(SYSTEM_SETTING_NAME)) {
      router.addRoute(MAIN_PATH_NAME, systemRoute)
      return next(to.fullPath)
    }

    // 移除系统路由
    if (!isAdmin.value && router.hasRoute(SYSTEM_SETTING_NAME)) {
      router.removeRoute(SYSTEM_SETTING_NAME)
      return next(to.fullPath)
    }

    next()
  })
}
