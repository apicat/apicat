import type { Router } from 'vue-router'
import { showLoading, hideLoading } from '@/hooks/useLoading'
import { ProjectRoutes } from '@/router/project.router'
import { useProjectStore } from '@/stores/project'
import { NOT_FOUND } from '../constant'

/**
 * 项目设置页面权限
 * @param router
 */
export default function initProjectNavFilter(router: Router) {
    router.beforeEach(async (to, from, next) => {
        const projectStore = useProjectStore()

        // 不在project模块内
        if (!ProjectRoutes.some((route) => route.name === to.name)) {
            return next()
        }

        // 当前停留在同一个项目
        if (projectStore.projectInfo && to.params.project_id === String(projectStore.projectInfo.id)) {
            return next()
        }

        let nextRoute = null

        showLoading()
        try {
            const projectInfo = await projectStore.getProjectDetail(parseInt(to.params.project_id as string))
            const authRouters = ProjectRoutes.filter((item) => item.meta.role.indexOf(projectInfo.authority) !== -1)
            const hasAuth = authRouters.some((route) => route.name === to.name)
            if (!hasAuth) {
                nextRoute = NOT_FOUND
            }
        } catch (error) {
            nextRoute = NOT_FOUND
        } finally {
            hideLoading()
        }

        return nextRoute ? next(nextRoute) : next()
    })
}
