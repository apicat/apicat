import useApi from '@/hooks/useApi'
import useProjectGroupStore from '@/store/projectGroup'
import { ProjectGroup, ProjectGroupSelectKey } from '@/typings'

export const useProjectGroups = () => {
  const selectedGroupRef = ref<ProjectGroupSelectKey>('all')
  const projectGroupStore = useProjectGroupStore()
  const [isLoading, getProjectGroupsApi] = useApi(projectGroupStore.getProjectGroups)

  onMounted(async () => await getProjectGroupsApi())

  return {
    selectedGroupRef,
    isLoading,
    getProjectGroups: getProjectGroupsApi,
  }
}
