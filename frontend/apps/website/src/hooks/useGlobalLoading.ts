import { useAppStoreWithOut } from '@/store/app'

export function useGlobalLoading(callback?: () => Promise<void> | void) {
  const appStore = useAppStoreWithOut()

  callback && onMounted(async () => {
    try {
      appStore.showGlobalLoading()
      await callback?.()
    }
    finally {
      appStore.hideGlobalLoading()
    }
  })

  return {
    showGlobalLoading: appStore.showGlobalLoading,
    hideGlobalLoading: appStore.hideGlobalLoading,
  }
}
