import { ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'

export const usePage = (onPageChange: any) => {
    const route = useRoute()
    const router = useRouter()

    const _page = parseInt(route.query.page as string, 10)

    const page = ref(isNaN(_page) ? 1 : _page)

    watch(
        () => page.value,
        () => {
            router.push({ name: router.currentRoute.value.name as string, query: { page: page.value } })
            onPageChange && onPageChange()
        }
    )

    return {
        page,
    }
}
