import useApi from '@/hooks/useApi'
import { ProjectInfo } from '@/typings'
import { getProjectList } from '@/api/project'

export const useProjectList = (searchParams?: Record<string, any>) => {
  const projectList = ref<ProjectInfo[]>([])
  const [isLoading, getProjectListApi] = useApi(getProjectList)

  // 加载项目列表
  const loadProjectList = async () => {
    try {
      projectList.value = await getProjectListApi(searchParams)
    } catch (error) {
      projectList.value = []
    }
  }

  onMounted(async () => await loadProjectList())
  onUnmounted(() => (projectList.value = []))

  return {
    isLoading,
    projects: projectList,
    projectList,

    getProjectList: loadProjectList,
  }
}
