import type { Router } from 'vue-router'
import { showLoading, hideLoading } from '@/hooks/useLoading'
import { useProjectStore } from '@/stores/project'
import { DOCUMENT_ROUTE_NAME, NOT_FOUND } from '../constant'

/**
 * 文档中项目权限拦截
 */
const initDocumentFilter = (route: Router) => {
    route.beforeEach(async (to, from, next) => {
        const projectStore = useProjectStore()

        // 不在文档模块内，忽略
        if (!to.matched.find((route) => route.name === DOCUMENT_ROUTE_NAME)) {
            return next()
        }

        const pid = parseInt(to.params.project_id as string, 10)

        // 同一个项目，无需再次获取权限信息
        if (projectStore.projectAuthInfo && projectStore.projectAuthInfo.id === pid) {
            return next()
        }

        try {
            showLoading()
            const projectAuthInfo = await projectStore.getProjectAuth(pid)

            // 无项目信息,404
            if (!projectAuthInfo || !projectAuthInfo.id) {
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

/**
 * 文档中项目详情获取
 */
export const initProjectDetail = (route: Router) => {
    route.beforeEach(async (to, from, next) => {
        const projectStore = useProjectStore()

        // 不在文档模块内，忽略
        if (!to.matched.find((route) => route.name === DOCUMENT_ROUTE_NAME)) {
            return next()
        }

        // console.log('2.获取项目详情')

        const pid = parseInt(to.params.project_id as string, 10)

        // 同一个项目，无需再次获取详情信息
        if (projectStore.projectInfo && projectStore.projectInfo.id === pid) {
            return next()
        }

        try {
            // showLoading()
            await projectStore.getProjectDetail(pid)
            // hideLoading()
        } catch (error) {
            hideLoading()
            return next(NOT_FOUND)
        }

        return next()
    })
}

export default initDocumentFilter
