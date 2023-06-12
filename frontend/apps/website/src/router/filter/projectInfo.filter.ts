import uesProjectStore from '@/store/project'
import { Router } from 'vue-router'
import { NOT_FOUND_PATH, NO_PERMISSION_PATH, PROJECT_DETAIL_PATH_NAME } from '../constant'
import { ProjectInfo } from '@/typings'
import { MemberAuthorityInProject } from '@/typings/member'

export const setupGetProjectInfoFilter = (router: Router) => {
  router.beforeEach(async (to, from, next) => {
    const projectStore = uesProjectStore()

    // filter project detail or edit page
    if (to.matched.find((item: any) => item.name === PROJECT_DETAIL_PATH_NAME)) {
      const project_id = to.params.project_id
      if (!projectStore.projectDetailInfo || projectStore.projectDetailInfo.id !== project_id) {
        try {
          const projectInfo: ProjectInfo = await projectStore.getProjectDetailInfo(project_id as string)

          if (!projectInfo) {
            return next(NOT_FOUND_PATH)
          }

          // 不在此项目中，无权限
          if (projectInfo.authority === MemberAuthorityInProject.NONE) {
            return next(NO_PERMISSION_PATH)
          }

          return next()
        } catch (error) {
          router.replace('/main')
        }
      }
    } else {
      // 非项目页，清空项目信息
      projectStore.clearCurrentProjectInfo()
    }

    next()
  })
}
