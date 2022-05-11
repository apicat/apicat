import { useApi } from '@/hooks/useApi'
import { onMounted, reactive, toRefs } from 'vue'
import { usePage } from '@/hooks/usePage'

export const useTable = (_api: any, options = { isLoaded: true, dataKey: 'data', totalKey: 'total' } as any) => {
    const { isLoaded = true, searchParam = {}, dataKey, totalKey } = options

    const [isLoading, api] = useApi(_api, { isShowMessage: false })

    const tableState = reactive({
        data: [],
        total: 0,
    })

    const { page } = usePage(() => getTableData())

    const getTableData = async () => {
        const res = await api({ ...searchParam, page: page.value })

        if (res && res.data) {
            tableState.data = res.data[dataKey] || []
            tableState.total = res.data[totalKey] || 1
        } else {
            tableState.data = []
            tableState.total = 0
        }
    }

    if (isLoaded) {
        onMounted(async () => {
            await getTableData()
        })
    }

    return {
        isLoading,
        currentPage: page,
        ...toRefs(tableState),
        getTableData,
    }
}

export default useTable
