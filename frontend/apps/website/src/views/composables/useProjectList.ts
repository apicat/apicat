import useApi from '@/hooks/useApi'
import { ProjectInfo } from '@/typings'
import { convertProjectCover, getProjectList } from '@/api/project'

export const useProjectList = (searchParams?: Record<string, any>) => {
  const projectList = ref<ProjectInfo[]>([])
  const [isLoading, getProjectListApi] = useApi(getProjectList)

  // 加载项目列表
  const loadProjectList = async () => {
    try {
      const projects = await getProjectListApi(searchParams)
      projectList.value = (projects || []).map((item: ProjectInfo) => convertProjectCover(item))
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
