export function useLoading() {
  const loadingRef = ref(0)

  return {
    loadingForGetter: () => loadingRef.value > 0,
    startLoading: () => {
      loadingRef.value++
    },
    endLoading: () => {
      loadingRef.value--
    },
  }
}
