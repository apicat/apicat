import type { Router } from 'vue-router'
import { getCollectionShareStatus } from '@/api/shareCollection'
import { useShareStore } from '@/store/share'
import { NOT_FOUND_PATH } from '../constant'

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
