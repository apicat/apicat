const title: Ref<string> = ref('')

export function useTitle(): Ref<string> {
  watch(title, () => {
    if (title.value && title.value !== 'ApiCat')
      document.title = `${title.value} - ApiCat`
    else
      document.title = 'ApiCat'
  })
  return title
}
