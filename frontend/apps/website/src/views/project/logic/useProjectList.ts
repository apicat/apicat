import useApi from '@/hooks/useApi'
import { convertProjectCover, getProjectList, toggleFollowProject } from '@/api/project'
import { ProjectInfo } from '@/typings'

export const useProjectList = (searchParams?: Record<string, any>) => {
  const projectList = ref<ProjectInfo[]>([])
  const [isLoading, getProjectListApi] = useApi(getProjectList)

  const handleFollowProject = async (project: ProjectInfo) => {
    try {
      await toggleFollowProject(project)
      project.is_followed = !project.is_followed
    } catch (error) {
      //
    }
  }

  onMounted(async () => {
    try {
      const projects = await getProjectListApi(searchParams)
      projectList.value = (projects || []).map((item: ProjectInfo) => convertProjectCover(item))
    } catch (error) {
      //
    }
  })
  onUnmounted(() => (projectList.value = []))

  return {
    isLoading,
    projectList,

    handleFollowProject,
  }
}
