import { ref } from 'vue'
import { useRoute } from 'vue-router'

export function usePage(pageSize?: number) {
  const route = useRoute()

  const queryPage = Number.parseInt(route.query.page as string, 10)
  const pageRef = ref(Number.isNaN(queryPage) ? 1 : queryPage)
  const pageSizeRef = ref(pageSize || 15)

  return {
    pageRef,
    pageSizeRef,
  }
}
