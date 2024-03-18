import { MAIN_PATH } from '@/router'

const contextKey = Symbol('historyGoback')
interface HistoryLayoutContext {
  goBack: () => Promise<void> | void
}

export function useHistoryLayoutProvide(): HistoryLayoutContext {
  const context = {
    goBack: () => {

    },
  }

  provide<HistoryLayoutContext>(contextKey, context)

  return context
}
export function useHistoryLayoutContext(): HistoryLayoutContext {
  return inject(contextKey)!
}
