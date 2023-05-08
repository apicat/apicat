import { LOGIN_PATH } from '@/router'
import { Router } from 'vue-router'
import { useUserStore } from '@/store/user'
import { storeToRefs } from 'pinia'

export const setupAuthFilter = (router: Router) => {
  router.beforeEach((to, from, next) => {
    const userStore = useUserStore()
    const { isLogin } = storeToRefs(userStore)

    if (to.meta.ignoreAuth) {
      return next()
    }

    if (!isLogin.value) {
      return next(LOGIN_PATH)
    }

    next()
  })
}
