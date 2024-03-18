import type { Router } from 'vue-router'
import { storeToRefs } from 'pinia'
import { LOGIN_NAME, LOGIN_PATH, MAIN_PATH, REGISTER_NAME, ROOT_PATH_NAME, redirectParam } from '@/router'
import { useUserStore } from '@/store/user'

// 自动重定向至主页的路由集合
const AutoRedirectToMainRoute = [LOGIN_NAME, REGISTER_NAME, ROOT_PATH_NAME]

function addRedirect(base: string, to: string) {
  const urlObj = new URL(base, location.origin)
  urlObj.searchParams.set(redirectParam, encodeURIComponent(to))
  return urlObj.toString().replace(location.origin, '')
}

export function popRedirect(url?: string) {
  if (!url)
    return MAIN_PATH
  const urlObj = new URL(url, location.origin)
  const redirect = urlObj.searchParams.get(redirectParam)
  return redirect ? decodeURIComponent(redirect) : MAIN_PATH
}

export function setupAuthFilter(router: Router) {
  router.beforeEach((to, from, next) => {
    const { isLogin: isLoginRef } = storeToRefs(useUserStore())
    const isLogin = unref(isLoginRef)

    if (isLogin && AutoRedirectToMainRoute.includes(to.name as string))
      return next(MAIN_PATH)

    if (to.meta.ignoreAuth)
      return next()

    if (!isLogin)
      return next(addRedirect(LOGIN_PATH, to.fullPath))

    next()
  })
}
