import { storeToRefs } from 'pinia'
import useApi from '@/hooks/useApi'
import { apiGetProjectList } from '@/api/project'
import { useTeamStore } from '@/store/team'

export function useProjectList(searchParams?: ProjectAPI.RequestProject) {
  const projectList = ref<ProjectAPI.ResponseProject[]>([])
  const { currentID } = storeToRefs(useTeamStore())
  const [isLoading, getProjectListApi] = useApi(apiGetProjectList)

  // 加载项目列表
  const loadProjectList = async () => {
    try {
      projectList.value = (await getProjectListApi(currentID.value, searchParams)) || []
    }
    catch (error) {
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
