import { getFollowedProjectList } from '@/api/project'
import useApi from '@/hooks/useApi'
import { ProjectInfo, SelectedKey } from '@/typings'

export const useFollowedProjectList = () => {
  const selectedRef = ref<SelectedKey>('all')
  const selectedHistory: SelectedKey[] = ['all']

  const followedProjects = ref<ProjectInfo[]>([])
  const [isLoading, getFollowedProjectListApi] = useApi(getFollowedProjectList)

  const activeClass = (key: SelectedKey) => (selectedRef.value === key ? 'active' : '')

  const goBackSelected = () => {
    const backSelected = selectedHistory.pop()
    if (backSelected) {
      selectedRef.value = backSelected
    }
  }

  const removeSelected = () => {
    selectedRef.value = null
  }

  onMounted(async () => {
    try {
      followedProjects.value = (await getFollowedProjectListApi()) || []
    } catch (error) {
      //
    }
  })

  return {
    isLoading,
    selectedRef,
    selectedHistory,
    followedProjects,
    activeClass,
    goBackSelected,
    removeSelected,
  }
}
