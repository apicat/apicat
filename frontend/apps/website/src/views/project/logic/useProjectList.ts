import useApi from '@/hooks/useApi'
import { convertProjectCover, getProjectList, toggleFollowProject } from '@/api/project'
import { ProjectInfo } from '@/typings'

export const useProjectList = () => {
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

  onBeforeMount(async () => {
    try {
      const projects = await getProjectListApi()
      projectList.value = (projects || []).map((item: ProjectInfo) => convertProjectCover(item))
    } catch (error) {
      //
    }
  })
  onBeforeUnmount(() => (projectList.value = []))

  return {
    isLoading,
    projectList,

    handleFollowProject,
  }
}
