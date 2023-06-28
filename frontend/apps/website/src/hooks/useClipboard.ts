import { useI18n } from 'vue-i18n'

export const useClipboard = (
  copyText: Ref<string> | string,
  elCopyText = '',
  elCopiedText = '',
  timeout = 1000
): { handleCopy: () => Promise<void>; elCopyTextRef: Ref<string> } => {
  const { t } = useI18n()

  elCopyText = elCopyText || t('app.common.copy')
  elCopiedText = (elCopiedText || t('app.tips.copyed')) + '!'

  const elCopyTextRef = ref(elCopyText)
  let timer: ReturnType<typeof setTimeout> | null = null

  const handleCopy = async (): Promise<void> => {
    if (!navigator.clipboard) {
      console.error('clipboard API is not supported')
      return
    }

    await navigator.clipboard.writeText(unref(copyText))
    if (!timeout) {
      return
    }

    elCopyTextRef.value = elCopiedText
    timer && clearTimeout(timer)
    timer = setTimeout(() => {
      elCopyTextRef.value = elCopyText
    }, timeout)
  }

  return {
    handleCopy,
    elCopyTextRef,
  }
}

export default useClipboard
