import type { Router } from 'vue-router'
import { useProjectStore } from '@/stores/project'
import { ITERATE_DOCUMENT_EDIT_NAME } from '../constant'
import { goNotFound } from '@/common/utils'
/**
 * 文档编辑权限拦截
 */
export default function initDocumentEditFilter(route: Router) {
    route.beforeEach(async (to, from, next) => {
        // 文档编辑拦截
        if (to.name === ITERATE_DOCUMENT_EDIT_NAME) {
            const projectStore = useProjectStore()
            // console.log('2.迭代 - 项目编辑权限拦截')
            if (projectStore.isReader) {
                goNotFound()
                return false
            }
        }
        return next()
    })
}
