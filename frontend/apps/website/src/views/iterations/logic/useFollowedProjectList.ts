import { getFollowedProjectList } from '@/api/project'
import useApi from '@/hooks/useApi'
import { ProjectInfo } from '@/typings'

export const useFollowedProjectList = () => {
  const selectedProjectKeyRef = ref<number | null>(null)
  const followedProjects = ref<ProjectInfo[]>([])

  const [isLoading, getFollowedProjectListApi] = useApi(getFollowedProjectList)

  // 设置当前选中的项目
  const setSelectedProjectKey = (key: number | null) => {
    selectedProjectKeyRef.value = key
  }

  onMounted(async () => {
    try {
      followedProjects.value = (await getFollowedProjectListApi()) || []
    } catch (error) {
      //
    }
  })

  return {
    selectedProjectKeyRef,
    isLoading,
    followedProjects,

    setSelectedProjectKey,
  }
}
