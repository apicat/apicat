import type { Router } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { LOGIN_PATH, MAIN_PATH, REGISTE_PATH } from '../constant'

/**
 * 权限拦截器，动态添加路由
 * @param router
 */
export default function initPermissionsFilter(router: Router) {
    router.beforeEach(async (to, from, next) => {
        const userStore = useUserStore()
        const isLogin = userStore.isLogin

        if (to.meta.ignoreAuth) {
            // 已登录,访问登录|注册页
            if (isLogin && (to.path === LOGIN_PATH || to.path === REGISTE_PATH)) {
                next(MAIN_PATH)
                return
            }
            next()
            return
        }

        if (!isLogin) {
            // 重定向至登录页
            const redirectData: any = {
                path: LOGIN_PATH,
                replace: true,
            }

            next(redirectData)

            return
        }

        next()
    })
}
