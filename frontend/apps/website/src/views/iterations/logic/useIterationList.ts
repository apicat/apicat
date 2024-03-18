import { useI18n } from 'vue-i18n'
import { storeToRefs } from 'pinia'
import type { IterationTableProps } from '../components/IterationTable.vue'
import { useTablev2 } from '@/hooks/useTable'
import { AsyncMsgBox } from '@/components/AsyncMessageBox'
import { useTeamStore } from '@/store/team'
import { apiDeleteIterationInfo, apiGetIterations } from '@/api/iteration/index'
import { ITERATION_DETAIL_PATH_NAME } from '@/router'

export function useIterationList(props: IterationTableProps) {
  const { t } = useI18n()
  const teamStore = useTeamStore()
  const { currentID } = storeToRefs(teamStore)
  const router = useRouter()

  const editableItreationIdRef = ref<number | string | null>(null)

  const queryParam: Record<string, any> = {
    projectID: props.projectId,
    teamID: currentID.value,
  }

  const { currentPage, getTableData, ...rest } = useTablev2(apiGetIterations, {
    pageSize: 10,
    isLoaded: false,
    addonArgs: [currentID.value],
  })

  const handleRemoveIteration = (id: string) => {
    AsyncMsgBox({
      title: t('app.iter.table.delete.title'),
      content: t('app.iter.table.delete.tip'),
      confirmButtonText: t('app.common.delete'),
      onOk: async () => {
        await apiDeleteIterationInfo(id)
        await getTableData(queryParam)
      },
    })
  }

  const handleRowClick = (iteration: IterationAPI.ResponseIteration) => {
    router.push({
      name: ITERATION_DETAIL_PATH_NAME,
      params: {
        iterationID: iteration.id,
      },
    })
  }

  // 项目切换时获取当前项目的迭代列表
  watch(
    () => props.projectId,
    async (projectID) => {
      queryParam.projectID = projectID || ''
      currentPage.value = 1
      await getTableData(queryParam)
    },
    {
      immediate: true,
    },
  )

  return {
    editableItreationIdRef,
    currentPage,
    ...rest,
    fetchIterationList: () => getTableData(queryParam),
    handleRemoveIteration,
    handleRowClick,
  }
}
