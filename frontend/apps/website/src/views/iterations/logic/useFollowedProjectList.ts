import { getFollowedProjectList } from '@/api/project'
import useApi from '@/hooks/useApi'
import { ProjectInfo, SelectedProjectKey } from '@/typings'

export const useFollowedProjectList = () => {
  const selectedProjectKeyRef = ref<SelectedProjectKey>(null)
  const followedProjects = ref<ProjectInfo[]>([])
  const [isLoading, getFollowedProjectListApi] = useApi(getFollowedProjectList)

  onMounted(async () => {
    try {
      followedProjects.value = (await getFollowedProjectListApi()) || []
    } catch (error) {
      //
    }
  })

  return {
    isLoading,
    selectedProjectKeyRef,
    followedProjects,
  }
}
