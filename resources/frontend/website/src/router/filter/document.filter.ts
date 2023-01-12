import type { Router } from 'vue-router'
import { useProjectStore } from '@/stores/project'
import { DOCUMENT_ROUTE_NAME } from '../constant'
import { goNotFound } from '@/common/utils'

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

        const project_id = parseInt(to.params.project_id as string, 10)

        // 同一个项目，无需再次获取权限信息
        if (projectStore.projectAuthInfo && projectStore.projectAuthInfo.id === project_id) {
            return next()
        }

        try {
            // showLoading()
            // console.log('1.获取项目状态')

            const projectAuthInfo = await projectStore.getProjectAuth(project_id)

            // 无项目信息,404
            if (!projectAuthInfo || !projectAuthInfo.id_public) {
                goNotFound()
                return false
            }
        } catch (error) {
            goNotFound()
            return false
        }

        return next()
    })
}

/**
 * 文档中项目详情获取
 */
export const initProjectDetailFilter = (route: Router) => {
    route.beforeEach(async (to, from, next) => {
        const projectStore = useProjectStore()

        // 不在文档模块内，忽略
        if (!to.matched.find((route) => route.name === DOCUMENT_ROUTE_NAME)) {
            return next()
        }

        // console.log('3.获取项目详情')

        const project_id = parseInt(to.params.project_id as string, 10)

        // 同一个项目，无需再次获取详情信息
        if (projectStore.projectInfo && projectStore.projectInfo.id === project_id) {
            return next()
        }

        try {
            await projectStore.getProjectDetail(project_id)
        } catch (error) {
            goNotFound()
            return false
        }

        return next()
    })
}

export default initDocumentFilter
