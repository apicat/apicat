import { useProjectStore } from '@/stores/project'
import { useIterateStore } from '@/stores/iterate'
import { useRoute } from 'vue-router'

const useIdPublicParam = () => {
    const { projectInfo } = useProjectStore()
    const { iterateInfo, isIterateRoute } = useIterateStore()

    const res = { isIterateRoute } as any

    res.projectId = projectInfo.id
    res.projectPublicId = projectInfo.id_public

    // 是否在迭代路由中
    if (isIterateRoute) {
        res.iterationPublicId = iterateInfo.id_public
        res.iterationId = iterateInfo.id
    }
    return res
}

/**
 * 根据当前路由获取对应的id_public
 */
export const getIdPublicByRouter = () => {
    const { isIterateRoute } = useIterateStore()
    const { params } = useRoute()
    return params[isIterateRoute ? 'iterate_id_public' : 'project_id_public']
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
