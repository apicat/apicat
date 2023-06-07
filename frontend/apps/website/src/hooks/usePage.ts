import { ref } from 'vue'
import { useRoute } from 'vue-router'

export const usePage = (pageSize?: number) => {
  const route = useRoute()

  const queryPage = parseInt(route.query.page as string, 10)
  const pageRef = ref(isNaN(queryPage) ? 1 : queryPage)
  const pageSizeRef = ref(pageSize || 15)

  return {
    pageRef,
    pageSizeRef,
  }
}
