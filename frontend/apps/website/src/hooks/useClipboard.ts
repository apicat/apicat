import { useI18n } from 'vue-i18n'

export const useClipboard = (
  copyText: Ref<string> | string,
  elCopyText = '',
  elCopiedText = '',
  timeout = 1000
): { handleCopy: () => Promise<void>; elCopyTextRef: Ref<string>; isCopied: Ref<boolean>; elCopiedText: string } => {
  const { t } = useI18n()
  const isClipboardApiSupported = useSupported(() => navigator && 'clipboard' in navigator)
  const isCopied: Ref<boolean> = ref(false)

  elCopyText = elCopyText || t('app.common.copy')
  elCopiedText = elCopiedText || t('app.tips.copyed')

  const elCopyTextRef = ref(elCopyText)
  let timer: ReturnType<typeof setTimeout> | null = null

  const handleCopy = async (): Promise<void> => {
    if (isClipboardApiSupported.value) {
      await navigator.clipboard.writeText(unref(copyText))
    } else {
      legacyCopy(unref(copyText))
    }

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

function legacyCopy(value: string) {
  const ta = document.createElement('textarea')
  ta.value = value ?? ''
  ta.style.position = 'absolute'
  ta.style.opacity = '0'
  document.body.appendChild(ta)
  ta.select()
  document.execCommand('copy')
  ta.remove()
}

export default useClipboard
