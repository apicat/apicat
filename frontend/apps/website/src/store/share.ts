import { pinia } from '@/plugins'
import { SharedDocumentInfo } from '@/typings'
import { defineStore } from 'pinia'

interface ShareState {
  sharedDocumentInfo: SharedDocumentInfo | null
}

export const uesShareStore = defineStore('share', {
  state: (): ShareState => ({
    sharedDocumentInfo: null,
  }),

  actions: {
    setDocumentShareInfo(info: SharedDocumentInfo) {
      this.sharedDocumentInfo = info
    },
    clearDocumentShareInfo() {
      this.sharedDocumentInfo = null
    },
  },
})

export default uesShareStore

export const uesShareStoreWithOut = () => uesShareStore(pinia)
