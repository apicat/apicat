import type { Router } from 'vue-router'
import { storeToRefs } from 'pinia'
import { useProjectStore } from '@/stores/project'
import { ITERATE_DOCUMENT_DETAIL_NAME } from '../constant'
import { useUserStore } from '@/stores/user'
import { unref } from 'vue'
import { goNotFound } from '@/common/utils'
/**
 * 文档详情(预览)权限拦截
 */

export default function initDocumentEditFilter(route: Router) {
    route.beforeEach(async (to, from, next) => {
        const { name } = to
        // 文档详情拦截
        if (name === ITERATE_DOCUMENT_DETAIL_NAME) {
            const projectStore = useProjectStore()
            const { projectAuthInfo, isPrivate: isPrivateRef } = storeToRefs(projectStore)
            const { has_shared, in_this } = unref(projectAuthInfo)
            const isPrivate = unref(isPrivateRef)
            const { isLogin: hasLogin } = useUserStore()

            // console.log(
            //     '2.迭代 - 项目详情权限拦截：',
            //     `\n\r是否分享:${has_shared ? '已分享' : '未分享'}`,
            //     `\n\r是否公开:${isPrivate ? '私有' : '公开'}`,
            //     `\n\r是否在项目中:${in_this ? '在' : '不在'}`
            // )

            // 公开项目
            if (!isPrivate) {
                next()
                return
            }

            // 已登录，在项目中
            if (hasLogin && in_this) {
                next()
                return
            }

            // 未公开 未登录 未分享 -> 404
            if (!hasLogin && isPrivate && !has_shared) {
                goNotFound()
                return false
            }

            // 未公开 未分享 已登录 不在项目中
            if (isPrivate && !has_shared && hasLogin && !in_this) {
                goNotFound()
                return false
            }

            // 未公开 未分享 已登录 在项目中
            if (isPrivate && !has_shared && hasLogin && in_this) {
                next()
                return
            }

            // 未公开 已分享 未输入密钥 -> 输秘钥
            // if (isPrivate && has_shared && !hasProjectSecretKey) {
            //     hideLoading()
            //     next({
            //         name: PREVIEW_PROJECT_SECRET,
            //         params: { id_public: params.id_public },
            //     })
            //     return
            // }
            //
            // // 未公开 已分享 已输入密钥|密钥未过期
            // if (isPrivate && has_shared && hasProjectSecretKey) {
            //     next()
            //     return
            // }

            goNotFound()
            return false
        }

        return next()
    })
}
