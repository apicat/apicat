import type { Router } from 'vue-router'
import { useProjectStore } from '@/stores/project'
import { DOCUMENT_EDIT_NAME, NOT_FOUND } from '../constant'
import { hideLoading } from '@/hooks/useLoading'
/**
 * 文档编辑权限拦截
 */
export default function initDocumentEditFilter(route: Router) {
    route.beforeEach(async (to, from, next) => {
        // 文档编辑拦截
        if (to.name === DOCUMENT_EDIT_NAME) {
            const projectStore = useProjectStore()
            if (projectStore.isReader) {
                hideLoading()
                return next(NOT_FOUND)
            }
        }
        return next()
    })
}
