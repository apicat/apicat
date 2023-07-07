import type { Router } from 'vue-router'
import { getCollectionShareStatus } from '@/api/shareCollection'
import { Cookies } from '@/commons'
import { useShareStore } from '@/store/share'
import { NOT_FOUND_PATH, PROJECT_SHARE_VALIDATION_NAME } from '../constant'
import { getDocumentShareDetailPath, getDocumentVerificationPath } from '../share'
import { getProjectDetailPath } from '../project.detail'

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
      const isExistSecretKey = !!(Cookies.get(Cookies.KEYS.SHARE_DOCUMENT + params.doc_public_id) || '')
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

        // 已经输入密钥，直接跳转文档详情
        if (name === 'share.document.verification' && isExistSecretKey) {
          return next(getDocumentShareDetailPath(params.doc_public_id as string))
        }

        // 未输入密钥访问该文档
        if (!isExistSecretKey && name === 'share.document.detail') {
          return next(getDocumentVerificationPath(params.doc_public_id as string))
        }

        return next()
      } catch (error) {
        return next(NOT_FOUND_PATH)
      }
    }

    const isExistSecretKeyForProject = !!(Cookies.get(Cookies.KEYS.SHARE_PROJECT + params.project_id) || '')
    // 项目密钥
    if (String(name).startsWith(PROJECT_SHARE_VALIDATION_NAME) && isExistSecretKeyForProject) {
      console.log(getProjectDetailPath(params.project_id as string))
      // return next(getProjectDetailPath(params.project_id as string))
    }

    next()
  })
}
