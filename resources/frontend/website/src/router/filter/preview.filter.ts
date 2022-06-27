import type { Router } from 'vue-router'
import { getDocumentStatus } from '@/api/preview'
import { showLoading, hideLoading } from '@/hooks/useLoading'
import { Storage, isNumber } from '@natosoft/shared'
import { PREVIEW_DOCUMENT, PREVIEW_DOCUMENT_SECRET, NOT_FOUND } from '../constant'
/**
 * 预览页面拦截
 * 已登录 | 未登录
 * @param router
 */
export default function initPreviewFilter(router: Router) {
    router.beforeEach(async (to, from, next) => {
        const { name, params } = to

        const hasDocumentSecretKey = !!(Storage.get(Storage.KEYS.SECRET_DOCUMENT_TOKEN + params.doc_id, true) || '')

        // 预览文档
        if (String(name).startsWith(PREVIEW_DOCUMENT)) {
            if (!isNumber(params.doc_id as string)) {
                return next(NOT_FOUND)
            }

            showLoading()

            const { data: hasShared } = await getDocumentStatus(params.doc_id)

            // 未被分享 | 不存在
            if (!hasShared) {
                hideLoading()
                return next(NOT_FOUND)
            }

            if (name === PREVIEW_DOCUMENT_SECRET) {
                hideLoading()
                return next()
            }

            if (!hasDocumentSecretKey) {
                hideLoading()
                return next({
                    name: PREVIEW_DOCUMENT_SECRET,
                    params: { doc_id: params.doc_id },
                })
            }

            return next()
        }

        next()
    })
}
