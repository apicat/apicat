import type { Router } from 'vue-router'
import { getCollectionShareStatus } from '@/api/shareCollection'
import { useShareStore } from '@/store/share'
import { NOT_FOUND_PATH, PROJECT_DETAIL_PATH_NAME } from '../constant'
import useProjectStore from '@/store/project'
import { useUserStore } from '@/store/user'
import { storeToRefs } from 'pinia'

/**
 * 分享页面拦截
 * 已登录 | 未登录
 * @param router
 */
export const setupShareDocumentFilter = (router: Router) => {
  router.beforeEach(async (to, from, next) => {
    const { name, params } = to
    // 预览文档
    if (String(name).startsWith('share.document')) {
      const shareStore = useShareStore()

      try {
        const sharedDocumentInfo = await getCollectionShareStatus(params.doc_public_id as string)

        // 未被分享 | 不存在
        if (!sharedDocumentInfo) {
          shareStore.clearDocumentShareInfo()
          return next(NOT_FOUND_PATH)
        }

        sharedDocumentInfo.doc_public_id = params.doc_public_id as string
        shareStore.setDocumentShareInfo(sharedDocumentInfo)

        return next()
      } catch (error) {
        return next(NOT_FOUND_PATH)
      }
    }

    next()
  })
}

export const setupShareProjectDetailFilter = (route: Router) => {
  route.beforeEach(async (to, from, next) => {
    // 项目详情拦截
    if (to.matched.find((route) => route.name === PROJECT_DETAIL_PATH_NAME)) {
      const projectStore = useProjectStore()
      const userStore = useUserStore()

      const { projectAuthInfo: projectAuthInfoRef, hasInputSecretKey: hasInputSecretKeyRef } = storeToRefs(projectStore)
      const { isLogin: isLoginRef } = storeToRefs(userStore)

      const projectAuthInfo = unref(projectAuthInfoRef)
      const hasInputSecretKey = unref(hasInputSecretKeyRef)
      const isLogin = unref(isLoginRef)

      // 无项目详情权限
      if (!projectAuthInfo) {
        return next(NOT_FOUND_PATH)
      }

      const { hasShared, isPrivate, inThisProject } = projectAuthInfo

      // 公开项目
      if (!isPrivate) {
        return next()
      }

      // 已登录，在项目中
      if (isLogin && inThisProject) {
        return next()
      }

      // 未公开 未登录 未分享 -> 404
      if (!isLogin && isPrivate && !hasShared) {
        return next(NOT_FOUND_PATH)
      }

      // 未公开 未分享 已登录 不在项目中
      if (isPrivate && !hasShared && isLogin && !inThisProject) {
        return next(NOT_FOUND_PATH)
      }

      // 未公开 未分享 已登录 在项目中
      if (isPrivate && !hasShared && isLogin && inThisProject) {
        return next()
      }

      // 未公开 已分享 未输入密钥 -> 输秘钥
      if (isPrivate && hasShared && !hasInputSecretKey) {
        projectStore.showProjectSecretLayer()
        return next()
      }

      // 未公开 已分享 已输入密钥|密钥未过期
      if (isPrivate && hasShared && hasInputSecretKey) {
        return next()
      }

      return next(NOT_FOUND_PATH)
    }

    return next()
  })
}
