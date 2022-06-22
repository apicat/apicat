import type { Router } from 'vue-router'
import { usePreviewStore } from '@/stores/preview'
import { useUserStore } from '@/stores/user'
import { showLoading, hideLoading } from '@/hooks/useLoading'
import { Storage } from '@natosoft/shared'
import { PREVIEW_DOCUMENT, PREVIEW_DOCUMENT_SECRET, PREVIEW_PROJECT, PREVIEW_PROJECT_SECRET, NOT_FOUND } from '../constant'
import { storeToRefs } from 'pinia'
/**
 * 预览页面拦截
 * 已登录 | 未登录
 * @param router
 */
export default function initPreviewFilter(router: Router) {
    router.beforeEach(async (to, from, next) => {
        const previewStore = usePreviewStore()
        const { isLogin: hasLogin } = useUserStore()
        const { projectInfo } = storeToRefs(previewStore)
        const { name, params } = to

        const hasProjectSecretKey = !!(Storage.get(Storage.KEYS.SECRET_PROJECT_TOKEN + params.project_id, true) || '')
        const hasDocumentSecretKey = !!(Storage.get(Storage.KEYS.SECRET_DOCUMENT_TOKEN + params.doc_id, true) || '')

        // 预览项目
        if (String(name).startsWith(PREVIEW_PROJECT)) {
            if (!projectInfo.value) {
                try {
                    showLoading()
                    const projectInfo = await previewStore.getProjectInfo(params.project_id)
                    // 无项目信息
                    if (!projectInfo || !projectInfo.id) {
                        hideLoading()
                        return next(NOT_FOUND)
                    }
                } catch (error) {
                    hideLoading()
                    return next(NOT_FOUND)
                }
            }

            const { has_shared, visibility, in_this } = projectInfo.value || {}

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
            if (!visibility && !hasLogin && !has_shared) {
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

            // 停留在秘钥页面无需二次跳转
            if (name === PREVIEW_PROJECT_SECRET) {
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

            next(NOT_FOUND)

            return
        }

        // 预览文档
        if (String(name).startsWith(PREVIEW_DOCUMENT)) {
            if (name === PREVIEW_DOCUMENT_SECRET) {
                next()
                return
            }

            if (!hasDocumentSecretKey) {
                next({
                    name: PREVIEW_DOCUMENT_SECRET,
                    params: { doc_id: params.doc_id },
                })
                return
            }

            showLoading()
            try {
                const documentInfo = await previewStore.getDocumentInfo(params.doc_id)

                // 无文档信息
                if (!documentInfo || !documentInfo.id) {
                    return next(NOT_FOUND)
                }

                return next()
            } catch (error) {
                //
            } finally {
                hideLoading()
            }

            next(NOT_FOUND)
            return
        }

        next()
    })
}
