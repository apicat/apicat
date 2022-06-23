import type { Router } from 'vue-router'
// import { showLoading, hideLoading } from '@/hooks/useLoading'
import { useProjectStore } from '@/stores/project'
import { DOCUMENT_DETAIL_NAME, NOT_FOUND, PREVIEW_PROJECT_SECRET } from '../constant'
import { useUserStore } from '@/stores/user'
import { Storage } from '@natosoft/shared'
/**
 * 文档详情(预览)权限拦截
 */

export default function initDocumentEditFilter(route: Router) {
    route.beforeEach(async (to, from, next) => {
        const { name, params } = to
        // 文档详情拦截
        if (name === DOCUMENT_DETAIL_NAME) {
            const projectStore = useProjectStore()
            const { isLogin: hasLogin } = useUserStore()
            const hasProjectSecretKey = !!(Storage.get(Storage.KEYS.SECRET_PROJECT_TOKEN + params.project_id, true) || '')
            const { has_shared, visibility, in_this } = projectStore.projectAuthInfo

            console.log(
                '项目权限拦截：',
                `\n\r是否分享:${has_shared ? '已分享' : '未分享'}`,
                `\n\r是否公开:${visibility ? '私有' : '公开'}`,
                `\n\r是否在项目中:${in_this ? '在' : '不在'}`
            )

            // 公开项目
            if (visibility) {
                next()
                return
            }

            // 已登录，在项目中
            if (hasLogin && in_this) {
                next()
                return
            }

            // 未公开 未登录 未分享 -> 404
            if (!hasLogin && !visibility && !has_shared) {
                next(NOT_FOUND)
                return
            }

            // 未公开 未分享 已登录 不在项目中
            if (!visibility && !has_shared && hasLogin && !in_this) {
                next(NOT_FOUND)
                return
            }

            // 未公开 未分享 已登录 在项目中
            if (!visibility && !has_shared && hasLogin && in_this) {
                next()
                return
            }

            // 未公开 已分享 未输入密钥 -> 输秘钥
            if (!visibility && has_shared && !hasProjectSecretKey) {
                next({
                    name: PREVIEW_PROJECT_SECRET,
                    params: { project_id: params.project_id },
                })
                return
            }

            // 未公开 已分享 已输入密钥|密钥未过期
            if (!visibility && has_shared && hasProjectSecretKey) {
                next()
                return
            }

            return next(NOT_FOUND)
        }

        return next()
    })
}
