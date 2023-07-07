import type { Router } from 'vue-router'
import { storeToRefs } from 'pinia'
import { useProjectStore } from '@/store/project'
import { useUserStore } from '@/store/user'
import useShareStore from '@/store/share'
import { NOT_FOUND_PATH, PROJECT_SHARE_VALIDATION_NAME } from '../constant'
import { DOCUMENT_DETAIL_NAME } from '../document'
import { Cookies } from '@/commons'

/**
 * 文档详情(预览)权限拦截
 */

export const setupDocumentDetailFilter = (route: Router) => {
  route.beforeEach(async (to, from, next) => {
    const { name, params } = to

    // 文档详情拦截
    if (name === DOCUMENT_DETAIL_NAME) {
      const project_id = params.project_id as string
      const projectStore = useProjectStore()
      const userStore = useUserStore()

      const { projectAuthInfo: projectAuthInfoRef } = storeToRefs(projectStore)
      const { isLogin: isLoginRef } = storeToRefs(userStore)

      const projectAuthInfo = unref(projectAuthInfoRef)
      const isLogin = unref(isLoginRef)
      const hasInputSecretKey = !!(Cookies.get(Cookies.KEYS.SHARE_PROJECT + project_id) || '')

      // 无项目详情权限
      if (!projectAuthInfo) {
        return next(NOT_FOUND_PATH)
      }

      const { hasShared, isPrivate, inThisProject } = projectAuthInfo

      console.log(
        '项目权限详情：',
        `\n\r是否输入密钥:${hasInputSecretKey ? '已输入密钥' : '未输入密钥'}`,
        `\n\r是否分享:${hasShared ? '已分享' : '未分享'}`,
        `\n\r是否公开:${isPrivate ? '私有' : '公开'}`,
        `\n\r是否在项目中:${inThisProject ? '在' : '不在'}`
      )

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
        return next({
          name: PROJECT_SHARE_VALIDATION_NAME,
          params: { project_id },
        })
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
