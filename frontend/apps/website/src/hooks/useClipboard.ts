import { useI18n } from 'vue-i18n'

export const useClipboard = (
  copyText: Ref<string> | string,
  elCopyText = '',
  elCopiedText = '',
  timeout = 1000
): { handleCopy: () => Promise<void>; elCopyTextRef: Ref<string>; isCopied: Ref<boolean>; elCopiedText: string } => {
  const { t } = useI18n()
  const isCopied: Ref<boolean> = ref(false)

  elCopyText = elCopyText || t('app.common.copy')
  elCopiedText = elCopiedText || t('app.tips.copyed')

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
    isCopied.value = true
    timer && clearTimeout(timer)
    timer = setTimeout(() => {
      isCopied.value = false
      elCopyTextRef.value = elCopyText
    }, timeout)
  }

  return {
    elCopyTextRef,
    elCopiedText,
    isCopied,
    handleCopy,
  }
}

export default useClipboard
