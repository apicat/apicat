import { sortProjectGroup } from '@/api/projectGroup'
import useApi from '@/hooks/useApi'
import useProjectGroupStore from '@/store/projectGroup'
import { ProjectGroup, ProjectGroupSelectKey } from '@/typings'
import { storeToRefs } from 'pinia'

export const useProjectGroups = () => {
  const selectedGroupRef = ref<ProjectGroupSelectKey>('all')
  const projectGroupStore = useProjectGroupStore()
  const { projectGroups } = storeToRefs(projectGroupStore)

  const [isLoading, getProjectGroupsApi] = useApi(projectGroupStore.getProjectGroups)

  onMounted(async () => await getProjectGroupsApi())

  const handleDeleteProjectGroup = async (group: ProjectGroup) => {
    console.log('handleDeleteProjectGroup', group)
  }

  const handleRenameProjectGroup = async (group: ProjectGroup) => {
    console.log('handleRenameProjectGroup', group)
  }

  const handleSortProjectGroup = async () => {
    await sortProjectGroup(projectGroups.value.map((item) => item.id!))
  }

  return {
    selectedGroupRef,
    isLoading,
    projectGroups,

    getProjectGroups: getProjectGroupsApi,
    handleDeleteProjectGroup,
    handleRenameProjectGroup,
    handleSortProjectGroup,
  }
}
