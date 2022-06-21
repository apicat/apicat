import type { Router } from 'vue-router'
import { useProjectStore } from '@/stores/project'
import { DOCUMENT_EDIT_NAME, NOT_FOUND } from '../constant'

/**
 * 文档编辑权限拦截
 */
export default function initDocumentEditFilter(route: Router) {
    route.beforeEach(async (to, from, next) => {
        // 文档编辑拦截
        if (to.name === DOCUMENT_EDIT_NAME) {
            const projectStore = useProjectStore()
            console.log('文档编辑拦截，当前权限是否为阅读者：', projectStore.isReader ? '是' : '否')
            if (projectStore.isReader) {
                return next(NOT_FOUND)
            }
        }
        return next()
    })
}
