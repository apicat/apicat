import type { Router } from 'vue-router'
import { useProjectStore } from '@/stores/project'
import { ITERATE_ROUTE_NAME } from '../constant'
import { useIterateStore } from '@/stores/iterate'
import {goNotFound} from "@/common/utils";

/**
 * 迭代模块拦截器，主要负责获取项目信息的获取,
 * 此处的id_public是与迭代ID一一对应的。
 */
const initIterateFilter = (route: Router) => {
    route.beforeEach(async (to, from, next) => {
        const projectStore = useProjectStore()
        const iterateStore = useIterateStore()

        iterateStore.isIterateRoute = false

        // 不在迭代模块内，忽略
        if (!to.matched.find((route) => route.name === ITERATE_ROUTE_NAME)) {
            return next()
        }

        iterateStore.isIterateRoute = true

        const iterate_id: any = to.params.iterate_id

        // 同一个迭代项目，无需再次获取权限信息
        if (iterateStore.iterateInfo && iterateStore.iterateInfo.id === iterate_id) {
            return next()
        }

        // 首先获取迭代所关联的项目
        try {

            const iterateInfo = await iterateStore.getIterateInfo(iterate_id)
            // console.log('1.迭代详情',iterateInfo)

            // 无迭代信息,404
            if (!iterateInfo || !iterateInfo.id) {
                goNotFound()
                return false
            }

            const projectInfo = await projectStore.getProjectDetailById({ iteration_id: iterate_id })
            // console.log('1.迭代 - 获取项目详情', projectInfo)

            // 无项目信息,404
            if (!projectInfo || !projectInfo.id) {
                goNotFound()
                return false
            }

            await projectStore.getProjectAuth(projectInfo.id)
            // console.log('2.迭代 - 获取项目权限', projectAuthInfo)
        } catch (error) {
            goNotFound()
            return false
        }

        return next()
    })
}

export default initIterateFilter
