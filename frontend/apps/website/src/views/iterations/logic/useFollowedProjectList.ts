import { getFollowedProjectList } from '@/api/project'
import useApi from '@/hooks/useApi'
import { ProjectInfo, SelectedKey } from '@/typings'

export const useFollowedProjectList = () => {
  const selectedRef = ref<SelectedKey>('all')
  let selectedHistory: number | ProjectInfo = 0

  const followedProjects = ref<ProjectInfo[]>([])
  const [isLoading, getFollowedProjectListApi] = useApi(getFollowedProjectList)

  const activeClass = (key: SelectedKey) => (selectedRef.value === key ? 'active' : '')

  const goBackSelected = () => {
    selectedRef.value = selectedHistory === 0 ? 'all' : (selectedHistory as ProjectInfo)
  }

  const removeSelected = () => {
    selectedRef.value = null
  }

  const setSelectedHistory = (info: ProjectInfo | number) => {
    selectedHistory = info
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
    setSelectedHistory,
  }
}
