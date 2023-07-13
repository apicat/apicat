import { MAIN_PATH, NO_PERMISSION_PATH, PROJECT_DETAIL_PATH_NAME } from '@/router'
import { Router } from 'vue-router'
import { storeToRefs } from 'pinia'
import useProjectStore from '@/store/project'

export const setupMemberPermissionFilter = (router: Router) => {
  router.beforeEach((to, from, next) => {
    const projectStore = useProjectStore()

    if (to.matched.find((item: any) => item.name === PROJECT_DETAIL_PATH_NAME) && !projectStore.isShowProjectSecretLayer) {
      const { projectDetailInfo } = storeToRefs(projectStore)

      if (!projectDetailInfo.value) {
        return next(MAIN_PATH)
      }

      if (to.meta.editableRoles && !(to.meta.editableRoles as string[]).includes(projectDetailInfo.value!.authority!)) {
        return next(NO_PERMISSION_PATH)
      }
    }

    next()
  })
}
