import { useUserStore } from '@/store/user'
import { isEmpty } from 'lodash-es'
import { Router } from 'vue-router'

export const setupGetUserInfoFilter = (router: Router) => {
  router.beforeEach(async (to, from, next) => {
    const { userInfo, getUserInfo } = useUserStore()

    // 检查用户信息
    if (!to.meta.ignoreAuth && (!userInfo || isEmpty(userInfo))) {
      await getUserInfo()
    }

    next()
  })
}
