import { Router } from 'vue-router'

export const setupAuthFilter = (router: Router) => {
  router.beforeEach((to, from, next) => {
    if (to.meta.ignoreAuth) {
      return next()
    }

    // if (!userStore.isLogin) {
    //   location.href = LOGIN_PATH
    //   return
    // }

    next()
  })
}
