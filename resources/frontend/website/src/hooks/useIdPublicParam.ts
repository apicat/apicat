import { useProjectStore } from '@/stores/project'
import { useIterateStore } from '@/stores/iterate'
import { useRoute } from 'vue-router'

const useIdPublicParam = () => {
    const { projectInfo } = useProjectStore()
    const { iterateInfo, isIterateRoute } = useIterateStore()

    const res = { isIterateRoute } as any

    res.projectId = projectInfo.id
    res.projectPublicId = projectInfo.id

    // 是否在迭代路由中
    if (isIterateRoute) {
        res.iterationPublicId = iterateInfo.id
        res.iterationId = iterateInfo.id
    }
    return res
}

/**
 * 根据当前路由获取对应的id
 */
export const getIdPublicByRouter = () => {
    const { isIterateRoute } = useIterateStore()
    const { params } = useRoute()
    return params[isIterateRoute ? 'iterate_id' : 'project_id']
}

export const generateProjectOrIterateParams = (params: any) => {
    const res: any = {}
    if (params.projectId) {
        res.project_id = params.projectId
    }

    if (params.iterationId) {
        res.iteration_id = params.iterationId
    }

    return res
}

export default useIdPublicParam
