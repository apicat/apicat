import type { RouteRecordRaw } from 'vue-router'
import {
  NOT_FOUND_PATH,
  NO_PERMISSION_PATH,
  PROJECT_COLLECTION_PATH,
  PROJECT_COLLECTION_PATH_NAME,
  PROJECT_DETAIL_PATH,
  PROJECT_DETAIL_PATH_NAME,
  PROJECT_RESPONSE_PATH,
  PROJECT_RESPONSE_PATH_NAME,
  PROJECT_SCHEMA_PATH,
  PROJECT_SCHEMA_PATH_NAME,
} from './constant'
import { useGlobalLoading } from '@/hooks/useGlobalLoading'
import useProjectStore from '@/store/project'
import { Authority, Visibility } from '@/commons/constant'
import { NotFoundError, UnauthorizedError } from '@/api/error'
import { useUserStore } from '@/store/user'
import { clearProjectSharedToken, getProjectSharedToken } from '@/api/shareToken'
import { useTitle } from '@/hooks/useTitle'

function isProjectPublic(shareStatus: ShareAPI.ResponseProjectShareStatus) {
  return shareStatus.visibility === Visibility.Public
}
function isInProject(shareStatus: ShareAPI.ResponseProjectShareStatus) {
  return shareStatus.permission !== Authority.None
}
function isProjectShared(shareStatus: ShareAPI.ResponseProjectShareStatus) {
  return shareStatus.hasShare
}

const title = useTitle()

export const projectDetailRoute: RouteRecordRaw = {
  name: PROJECT_DETAIL_PATH_NAME,
  path: PROJECT_DETAIL_PATH,
  meta: { ignoreAuth: true },
  redirect: { name: PROJECT_COLLECTION_PATH_NAME },
  beforeEnter: async (to, _, next) => {
    const projectID = to.params.project_id as string
    const { showGlobalLoading, hideGlobalLoading } = useGlobalLoading()
    const projectStore = useProjectStore()
    const haveShareToken = getProjectSharedToken(projectID)

    showGlobalLoading()
    try {
      // before
      const isLogin = useUserStore().isLogin

      const shareStatus = await projectStore.getProjejctAuthInfo(projectID)
      let GetProjectInfoNeeded = false

      // jump
      //   项目不存在 -> x
      if (!shareStatus) {
        return next(NOT_FOUND_PATH)
      }
      //   项目公开 -> o
      else if (isProjectPublic(shareStatus)) {
        GetProjectInfoNeeded = true
      }
      //   是项目成员且已登陆 -> o
      else if (isLogin && isInProject(shareStatus)) {
        GetProjectInfoNeeded = true
      }
      //   项目未被分享 -> x
      else if (!isProjectShared(shareStatus)) {
        return next(NO_PERMISSION_PATH)
      }
      //   项目被分享 -> o
      else if (isProjectShared(shareStatus)) {
        if (haveShareToken)
          GetProjectInfoNeeded = true
      }
      else {
        return next(NOT_FOUND_PATH)
      }

      // 是否显示密钥输入层
      projectStore.isShowProjectSecretLayer = !GetProjectInfoNeeded

      // get info
      if (GetProjectInfoNeeded) {
        await projectStore.getProjectInfoById(projectID)
        title.value = projectStore.project!.title!
      }
    }
    catch (error) {
      // 401 - 清除share token
      if (error instanceof UnauthorizedError) {
        clearProjectSharedToken(projectID)
        projectStore.isShowProjectSecretLayer = true
      }
      // 404 - NotFoundError
      else if (error instanceof NotFoundError) {
        return next(NOT_FOUND_PATH)
      }
      // default - show error page
      else { return next(NOT_FOUND_PATH) }
    }
    finally {
      hideGlobalLoading()
    }

    next()
  },

  component: async () => import('@/layouts/ProjectDetailLayout/ProjectDetailLayout.vue'),
  children: [
    {
      name: PROJECT_COLLECTION_PATH_NAME,
      path: PROJECT_COLLECTION_PATH,
      component: () => import('@/views/collection/CollectionPage.vue'),
      props: true,
    },
    {
      name: PROJECT_SCHEMA_PATH_NAME,
      path: PROJECT_SCHEMA_PATH,
      component: () => import('@/views/schema/SchemaPage.vue'),
      props: true,
    },
    {
      name: PROJECT_RESPONSE_PATH_NAME,
      path: PROJECT_RESPONSE_PATH,
      component: () => import('@/views/response/ResponsePage.vue'),
      props: true,
    },
  ],
}
