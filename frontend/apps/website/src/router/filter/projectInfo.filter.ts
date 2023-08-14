import useProjectStore from '@/store/project'
import type { Router } from 'vue-router'
import { MAIN_PATH, NOT_FOUND_PATH, PROJECT_DETAIL_PATH_NAME, ITERATION_DETAIL_PATH_NAME, ITERATION_LIST_ROOT_PATH } from '../constant'
import { ProjectInfo } from '@/typings'
import { useIterationStore } from '@/store/iteration'

export const setupGetProjectAuthInfoFilter = (router: Router) => {
  router.beforeEach(async (to, from, next) => {
    const projectStore = useProjectStore()

    if (to.matched.find((route) => route.name === PROJECT_DETAIL_PATH_NAME)) {
      // 同一个项目，无需再次获取权限信息
      if (projectStore.projectAuthInfo && projectStore.projectAuthInfo.project_id === to.params.project_id) {
        return next()
      }

      try {
        await projectStore.getProjectAuthInfo(to.params.project_id as string)
        return next()
      } catch (error) {
        return next(MAIN_PATH)
      }
    }

    next()
  })
}

export const setupGetProjectInfoFilter = (router: Router) => {
  router.beforeEach(async (to, from, next) => {
    const projectStore = useProjectStore()

    if (to.matched.find((item: any) => item.name === PROJECT_DETAIL_PATH_NAME) && !projectStore.isShowProjectSecretLayer) {
      const project_id = to.params.project_id
      if (!projectStore.projectDetailInfo || projectStore.projectDetailInfo.id !== project_id) {
        try {
          const projectInfo: ProjectInfo = await projectStore.getProjectDetailInfo(project_id as string)

          if (!projectInfo) {
            return next(NOT_FOUND_PATH)
          }

          return next()
        } catch (error) {
          return next(MAIN_PATH)
        }
      }
    }

    next()
  })
}

/**
 * 通过迭代信息获取项目信息&权限信息
 * @param router
 */
export const setupGetProjectInfoByIterationFilter = (router: Router) => {
  router.beforeEach(async (to, from, next) => {
    const iterationStore = useIterationStore()
    const projectStore = useProjectStore()

    if (to.matched.find((item) => item.name === ITERATION_DETAIL_PATH_NAME)) {
      const iteration_id = to.params.iteration_id

      if (!iterationStore.iterationInfo || iterationStore.iterationInfo.id !== iteration_id) {
        try {
          const iterationInfo = await iterationStore.getIterationInfo(iteration_id as string)
          await projectStore.getProjectDetailInfo(iterationInfo.project_id)
          await projectStore.getProjectAuthInfo(iterationInfo.project_id)

          if (!iterationStore.iterationInfo || !projectStore.projectDetailInfo || !projectStore.projectAuthInfo) {
            return next(NOT_FOUND_PATH)
          }

          return next()
        } catch (error) {
          return next(ITERATION_LIST_ROOT_PATH)
        }
      }
    }

    next()
  })
}
