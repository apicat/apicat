import { storeToRefs } from 'pinia'
import { apiGetMyFollowedProjectList } from '@/api/project'
import useApi from '@/hooks/useApi'
import { useTeamStore } from '@/store/team'
import { Storage } from '@/commons'
import type { IterationTreeEmits } from '@/views/iterations/components/IterationTree.vue'

export function useFollowedProjectList(emit: IterationTreeEmits) {
  const { currentID } = storeToRefs(useTeamStore())

  const statusKey = 'selectedIterationProject'
  const getStatus = () => Storage.get(statusKey) || 'all'
  const setStatus = (key: IterationSelectedKey) => {
    if (key === 'create')
      key = 'all'
    Storage.set(statusKey, key)
    return key
  }
  const selectedRef = ref<IterationSelectedKey>(getStatus())
  let selectedHistory: IterationSelectedKey = selectedRef.value

  const followedProjects = ref<ProjectAPI.ResponseProject[]>([])
  const [isLoading, getFollowedProjectListApi] = useApi(apiGetMyFollowedProjectList)

  const activeClass = (key: IterationSelectedKey) => (selectedRef.value === key ? 'active' : '')

  const goBackSelected = () => (selectedRef.value = selectedHistory)

  const removeSelected = () => (selectedRef.value = 'all')

  const goSelectedAll = () => {
    selectedHistory = 'all'
    goBackSelected()
  }

  // 保存打开状态
  const setSelectedHistory = (info: IterationSelectedKey) => {
    info = setStatus(info)
    selectedHistory = info
  }

  onMounted(async () => {
    try {
      followedProjects.value = (await getFollowedProjectListApi(currentID.value)) || []
      if (selectedRef.value === 'all' || selectedRef.value === 'create')
        return
      let selectedProject: ProjectAPI.ResponseProject | undefined
      for (let i = 0; i < followedProjects.value.length; i++) {
        if (followedProjects.value[i].id === selectedRef.value) {
          selectedProject = followedProjects.value[i]
          break
        }
      }
      if (!selectedProject)
        selectedRef.value = 'all'
      else emit('click', selectedProject)
    }
    catch (error) {
      //
    }
  })

  return {
    isLoading,
    selectedRef,
    selectedHistory,
    followedProjects,
    activeClass,
    goBackSelected,
    goSelectedAll,
    removeSelected,
    setSelectedHistory,
  }
}
