import useTable from '@/hooks/useTable'
import { getIterationList, deleteIteration } from '@/api/iteration'
import { Iteration } from '@/typings'
import { AsyncMsgBox } from '@/components/AsyncMessageBox'
import { useI18n } from 'vue-i18n'
import { getIterationDetailPath } from '@/router/iteration.detail'

export const useIterationList = (projectIdRef: Ref<number | string | null>) => {
  const { t } = useI18n()
  const router = useRouter()

  const editableItreationIdRef = ref<number | string | null>(null)

  const queryParam: Record<string, any> = {
    project_id: projectIdRef.value,
  }

  const { currentPage, getTableData, ...rest } = useTable<Iteration>(getIterationList, {
    searchParam: queryParam,
    isLoaded: false,
    dataKey: 'iterations',
  })

  const handleRemoveIteration = (iteration: Iteration) => {
    AsyncMsgBox({
      title: t('app.common.deleteTip'),
      content: '确定删除该迭代吗?',
      onOk: async () => {
        await deleteIteration({ iteration_id: iteration.id })
        await getTableData()
      },
    })
  }

  const handleRowClick = (iteration: Iteration) => {
    router.push(getIterationDetailPath(iteration.id))
  }

  // 项目切换时获取当前项目的迭代列表
  watch(
    projectIdRef,
    async () => {
      queryParam.project_id = projectIdRef.value || ''
      currentPage.value = 1
      await getTableData()
    },
    {
      immediate: true,
    }
  )

  return {
    editableItreationIdRef,
    currentPage,
    ...rest,
    fetchIterationList: getTableData,
    handleRemoveIteration,
    handleRowClick,
  }
}
