import type { InjectionKey } from 'vue'
import { inject, provide } from 'vue'
import type { PageModeCtx } from '@/views/composables/usePageMode'
import { usePageMode } from '@/views/composables/usePageMode'

interface PagesMode {
  collection: PageModeCtx
  schema: PageModeCtx
  response: PageModeCtx
}

export const pagesModeKey: InjectionKey<PagesMode> = Symbol('pagesModeKey')

export function providePagesMode() {
  provide(pagesModeKey, {
    collection: usePageMode(),
    schema: usePageMode(),
    response: usePageMode(),
  })
}

export function injectPagesMode(page?: keyof PagesMode) {
  return page ? inject(pagesModeKey)![page] : inject(pagesModeKey)
}
