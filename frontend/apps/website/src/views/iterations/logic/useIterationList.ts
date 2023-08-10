import useTable from '@/hooks/useTable'
import { getIterationList } from '@/api/iteration'
import { Iteration } from '@/typings'

export type SelectedProjectKey = number | string | null

export const useIterationList = () => {
  const selectedProjectKeyRef = ref<SelectedProjectKey>(null)

  const queryParam: Record<string, any> = {
    project_id: selectedProjectKeyRef.value,
  }

  const { currentPage, getTableData, ...rest } = useTable<Iteration>(getIterationList, {
    searchParam: queryParam,
    isLoaded: false,
    dataKey: 'iterations',
  })

  // 项目切换时获取当前项目的迭代列表
  watch(
    selectedProjectKeyRef,
    async () => {
      queryParam.project_id = selectedProjectKeyRef.value || ''
      currentPage.value = 1
      await getTableData()
    },
    {
      immediate: true,
    }
  )

  return {
    selectedProjectKeyRef,
    currentPage,
    ...rest,
    fetchIterationList: getTableData,
  }
}
