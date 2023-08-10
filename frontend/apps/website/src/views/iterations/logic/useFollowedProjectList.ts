import { getFollowedProjectList } from '@/api/project'
import useApi from '@/hooks/useApi'
import { ProjectInfo } from '@/typings'

export const useFollowedProjectList = () => {
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
    followedProjects,
  }
}
