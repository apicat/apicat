import type { Router } from 'vue-router'
import { usePreviewStore } from '@/stores/preview'
import { showLoading, hideLoading } from '@/hooks/useLoading'
import { Storage } from '@natosoft/shared'
import { PREVIEW_DOCUMENT, PREVIEW_DOCUMENT_SECRET, NOT_FOUND } from '../constant'
/**
 * 预览页面拦截
 * 已登录 | 未登录
 * @param router
 */
export default function initPreviewFilter(router: Router) {
    router.beforeEach(async (to, from, next) => {
        const previewStore = usePreviewStore()
        const { name, params } = to

        const hasDocumentSecretKey = !!(Storage.get(Storage.KEYS.SECRET_DOCUMENT_TOKEN + params.doc_id, true) || '')

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
