import { onMounted, reactive, toRefs, watch } from 'vue'
import type { WatchStopHandle } from 'vue'
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

export function useTable<T extends Record<string, any>, R>(_api: (args: T) => Promise<R>, options: UseTableOptions) {
  const { isLoaded = true, searchParam = {}, dataKey = 'items', totalKey = 'count', pageSize = 15, transform } = options

  const [isLoading, api] = useApi(_api, { isShowMessage: false })

  const tableState: {
    data: R[]
    total: number
  } = reactive({
    data: [],
    total: 0,
  })

  const { pageRef, pageSizeRef } = usePage(pageSize)

  const getTableData = async () => {
    try {
      const _searchParam = Object.keys(searchParam || {}).reduce((acc, key) => {
        if (searchParam[key])
          acc[key] = unref(searchParam[key])
        return acc
      }, {} as Record<string, any>)

      const data = await (api as any)({
        ..._searchParam,
        page: pageRef.value,
        page_size: pageSizeRef.value,
      })
      if (data) {
        tableState.data = (data[dataKey] || []).map((item: any) => (isFunction(transform) ? transform(item) : item))
        tableState.total = data[totalKey] || 0
      }
      else {
        tableState.data = []
        tableState.total = 0
      }
    }
    catch (error) {
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

interface UseTableOptionsv2<T, R> {
  isLoaded?: boolean
  dataProcess?: (data: R) => R
  pageSize: number
  args: T
  addonArgs: any[]
  doWatch: true
  dataKey: GlobalAPI.ResponseTableDataKey
  totalKey: string
}

export function useTablev2<T extends GlobalAPI.RequestTable, R, P extends any[]>(
  _api: (args: Partial<T>, ...addon: P) => Promise<GlobalAPI.ResponseTable<R[]>>,
  options: Partial<UseTableOptionsv2<T, R>>,
) {
  const {
    isLoaded = true,
    dataProcess,
    pageSize,
    args = <T>{},
    addonArgs = <any>[],
    doWatch = true,
    dataKey = 'items',
    totalKey = 'count',
  } = options || {}
  const [isLoading, api] = useApi(_api, { isShowMessage: false })
  const data = ref<R[]>([])
  const tableState = reactive<{
    totalPage: number
    total: number
    pageSize: number
    currentPage: number
  }>({
    totalPage: 0,
    total: 0,
    pageSize: pageSize || 10,
    currentPage: 0,
  })

  const { pageRef, pageSizeRef } = usePage(pageSize)

  const getTableData = async (p: Partial<T> | undefined, ...a: P | []) => {
    const d = p || args || ({} as Partial<T>)
    // const addon = a || addonArgs || <any>[]
    const addon = ((): P => {
      if (!addonArgs) {
        return a as P
      }
      else {
        if (a.length > 0)
          return a as P
        else
          return addonArgs as P
      }
    })()

    try {
      const res = await api(
        {
          ...d,
          page: pageRef.value,
          pageSize: pageSizeRef.value,
        },
        ...(addon as P),
      )
      if (!res) {
        data.value = []
        return
      }

      tableState.total = (res as any)[totalKey] || 0
      data.value = ((res as any)[dataKey] || [])
        .filter((val: R) => val)
        .map((item: R) => (isFunction(dataProcess) ? dataProcess(item) : item))
    }
    catch (error) {
      data.value = []
      tableState.total = 0
    }
  }

  let stopWatch: WatchStopHandle | undefined
  function startWatch() {
    if (!stopWatch)
      stopWatch = watch([pageRef, pageSizeRef], () => getTableData(args, ...(addonArgs as P)))
  }

  function unWatch() {
    if (stopWatch) {
      stopWatch()
      stopWatch = undefined
    }
  }

  if (doWatch) {
    onMounted(() => startWatch())
    onUnmounted(() => unWatch())
  }

  isLoaded && onMounted(async () => await getTableData(undefined))

  return {
    ...toRefs(tableState),
    currentPage: pageRef,
    pageSize: pageSizeRef,
    isLoading,
    data,

    getTableData,
    refreshData: () => getTableData(undefined),
    startWatch,
    unWatch,
  }
}
