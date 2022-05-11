import type { Router } from 'vue-router'
import { showLoading, hideLoading } from '@/hooks/useLoading'
import { useProjectStore } from '@/stores/project'
import { DOCUMENT_ROUTE_NAME, NOT_FOUND } from '../constant'

/**
 * 文档编辑权限拦截
 */

export default function initDocumentEditFilter(route: Router) {
    route.beforeEach(async (to, from, next) => {
        const projectStore = useProjectStore()

        // 不在document模块内,清除一缓存的项目信息
        if (!to.matched.find((route) => route.name === DOCUMENT_ROUTE_NAME)) {
            return next()
        }

        const pid = parseInt(to.params.project_id as string, 10)

        if (projectStore.projectInfo && projectStore.projectInfo.id === pid) {
            return next()
        }

        try {
            showLoading()
            await projectStore.getProjectDetail(pid)
            if (projectStore.isReader) {
                hideLoading()
                return next(NOT_FOUND)
            }
        } catch (error) {
            hideLoading()
            return next(NOT_FOUND)
        }

        return next()
    })
}
