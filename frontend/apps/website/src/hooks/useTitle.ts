import { useTitle as _useTitle } from '@vueuse/core'

const title: Ref<string> = ref('')

export const useTitle = () => {
  _useTitle(
    computed(() => {
      return !title?.value ? 'ApiCat' : `${title?.value} - ApiCat`
    })
  )

  return title
}
