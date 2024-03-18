import { useI18n } from 'vue-i18n'

export function useClipboard(
  copyText: Ref<string> | string,
  _elCopyText: ComputedRef<string> | string = '',
  _elCopiedText: ComputedRef<string> | string = '',
  timeout = 1000,
): {
    handleCopy: () => Promise<void>
    elCopyTextRef: Ref<string>
    isCopied: Ref<boolean>
    elCopiedText: ComputedRef<string>
  } {
  const { t } = useI18n()
  const isClipboardApiSupported = useSupported(() => navigator && 'clipboard' in navigator)
  const isCopied: Ref<boolean> = ref(false)

  const elCopyText = computed(() => {
    if (_elCopyText) {
      if (typeof _elCopyText === 'string')
        return _elCopyText

      else
        return _elCopyText.value
    }

    return t('app.common.copy')
  })
  const elCopiedText = computed(() => {
    if (_elCopiedText) {
      if (typeof _elCopiedText === 'string')
        return _elCopiedText

      else
        return _elCopiedText.value
    }
    return t('app.common.copyed')
  })

  const elCopyTextRef = ref(elCopyText.value)
  let timer: ReturnType<typeof setTimeout> | null = null

  const handleCopy = async (): Promise<void> => {
    if (isClipboardApiSupported.value)
      await navigator.clipboard.writeText(unref(copyText))
    else legacyCopy(unref(copyText))

    if (!timeout)
      return

    elCopyTextRef.value = elCopiedText.value
    isCopied.value = true
    timer && clearTimeout(timer)
    timer = setTimeout(() => {
      isCopied.value = false
      elCopyTextRef.value = elCopyText.value
    }, timeout)
  }

  watch(elCopyText, () => {
    elCopyTextRef.value = elCopyText.value
  })

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
