import { Cookies } from '@/commons'
import { pinia } from '@/plugins'
import { SharedDocumentInfo } from '@/typings'
import { defineStore } from 'pinia'

interface ShareState {
  sharedDocumentInfo: SharedDocumentInfo | null
}

export const useShareStore = defineStore('share', {
  state: (): ShareState => ({
    sharedDocumentInfo: null,
  }),

  getters: {
    token() {
      const currentRouteMatched = this.$router.currentRoute.value.matched
      const params = this.$router.currentRoute.value.params
      const { doc_public_id, project_id } = params as Record<string, string>

      // 预览分享的文档
      if (currentRouteMatched.find((route) => route.name === 'share.document')) {
        return Cookies.get(Cookies.KEYS.SHARE_DOCUMENT + (doc_public_id || ''))
      }

      // 预览分享的项目
      if (currentRouteMatched.find((route) => route.name === 'project.detail')) {
        return Cookies.get(Cookies.KEYS.SHARE_PROJECT + (project_id || ''))
      }

      return null
    },
  },

  actions: {
    setDocumentShareInfo(info: SharedDocumentInfo) {
      this.sharedDocumentInfo = info
    },
    clearDocumentShareInfo() {
      this.sharedDocumentInfo = null
    },
  },
})

export default useShareStore

export const useShareStoreWithOut = () => useShareStore(pinia)

export const setShareTokenToParams = (params: Record<string, any>): Record<string, any> => {
  const shareStore = useShareStoreWithOut()
  if (shareStore.token) {
    params.token = shareStore.token
  }
  return params
}
