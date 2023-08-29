type PageMode = 'list' | 'form'

export const usePageMode = (initMode: PageMode = 'list') => {
  const modeRef: Ref<PageMode> = ref(initMode)

  const switchMode = (m: PageMode) => {
    modeRef.value = m
  }

  const toggleMode = () => {
    modeRef.value = modeRef.value === 'list' ? 'form' : 'list'
  }

  const isListMode = computed(() => modeRef.value === 'list')

  const isFormMode = computed(() => modeRef.value === 'form')

  return {
    switchMode,
    toggleMode,

    isListMode,
    isFormMode,
  }
}
