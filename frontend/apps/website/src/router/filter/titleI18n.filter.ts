import type { Router } from 'vue-router'
import { useTitle } from '@/hooks/useTitle'
import useLocaleStore from '@/store/locale'

export function setupTitleI18n(router: Router) {
  const title = useTitle()
  router.beforeEach((to, from, next) => {
    const { t } = useLocaleStore()
    title.value = to.meta.title ? t(to.meta.title as string) : 'ApiCat'
    next()
  })
}
