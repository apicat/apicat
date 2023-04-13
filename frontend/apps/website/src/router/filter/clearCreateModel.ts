import uesAppStore from '@/store/app'
import { Router } from 'vue-router'

export const setupClearCreateModelFilter = (router: Router) => {
  // 清空当前的创建模式
  router.beforeEach(() => {
    const appStore = uesAppStore()
    appStore.setCreateMode(null)
  })
}
