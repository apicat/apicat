type PageMode = 'list' | 'form'

export function usePageMode(initMode: PageMode = 'list') {
  const modeRef: Ref<PageMode> = ref(initMode)

  const switchMode = (m: PageMode) => {
    modeRef.value = m
  }

  const toggleMode = () => {
    modeRef.value = modeRef.value === 'list' ? 'form' : 'list'
  }

  const isListMode = computed(() => modeRef.value === 'list')

  const isFormMode = computed(() => modeRef.value === 'form')

  const switchToWriteMode = () => switchMode('form')

  const switchToReadMode = () => switchMode('list')

  return {
    switchMode,
    toggleMode,
    switchToWriteMode,
    switchToReadMode,

    isListMode,
    isFormMode,

    isReadMode: isListMode,
    isWriteMode: isFormMode,
    readonly: isListMode,
  }
}

export interface PageModeCtx {
  switchMode: (m: PageMode) => void
  toggleMode: () => void
  switchToWriteMode: () => void
  switchToReadMode: () => void
  isListMode: ComputedRef<boolean>
  isFormMode: ComputedRef<boolean>
  isReadMode: ComputedRef<boolean>
  isWriteMode: ComputedRef<boolean>
  readonly: ComputedRef<boolean>
}
