import { onMounted, reactive, toRefs, watch } from 'vue'
import { isFunction } from 'lodash-es'
import { useApi } from '@/hooks/useApi'
import { usePage } from '@/hooks/usePage'

interface UseTableOptions {
  isLoaded?: boolean
  searchParam?: Record<string, any>
  dataKey?: string
  totalKey?: string
  pageSize?: number
  transform?: (item: any) => any
}

export const useTable = (_api: any, options: UseTableOptions) => {
  const { isLoaded = true, searchParam = {}, dataKey = 'records', totalKey = 'total', pageSize = 15, transform } = options

  const [isLoading, api] = useApi(_api, { isShowMessage: false })

  const tableState = reactive({
    data: [],
    total: 0,
  })

  const { pageRef, pageSizeRef } = usePage(pageSize)

  const getTableData = async () => {
    const data = await api({ ...searchParam, page: pageRef.value, page_size: pageSizeRef.value })
    if (data) {
      tableState.data = (data[dataKey] || []).map((item: any) => (isFunction(transform) ? transform(item) : item))
      tableState.total = data[totalKey] || 1
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

  watch([pageRef, pageSizeRef], () => getTableData())

  return {
    isLoading,
    currentPage: pageRef,
    pageSize: pageSizeRef,
    ...toRefs(tableState),

    getTableData,
  }
}

export default useTable
