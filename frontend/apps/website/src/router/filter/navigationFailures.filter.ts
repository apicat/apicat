import type { Router } from 'vue-router'
import { ElMessage } from 'element-plus'
import useLocaleStore from '@/store/locale'

export function setupNavigationFailureFilter(router: Router) {
  router.onError((e) => {
    const { t } = useLocaleStore()
    ElMessage.error(`${t('app.common.resourceLoadErrorTip')}: ${e}`)
  })
}
