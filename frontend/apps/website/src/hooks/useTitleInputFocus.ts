export function useTitleInputFocus() {
  const inputRef = ref<HTMLInputElement>()
  return {
    inputRef,
    focus() {
      nextTick(() => {
        inputRef.value?.focus()
      })
    },
  }
}
