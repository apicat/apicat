import { storeToRefs } from 'pinia'
import type SelectProjectGroup from '../components/SelectProjectGroup.vue'
import {
  apiFollowProject,
  apiGetMyFollowedProjectList,
  apiGetMyProjectList,
  apiGetProjectList,
  apiGetProjectListByGroupId,
  apiUnfollowProject,
} from '@/api/project'
import useProjectGroupStore from '@/store/projectGroup'
import { useTeamStore } from '@/store/team'
import { PROJECT_DETAIL_PATH_NAME } from '@/router'
import { useGlobalLoading } from '@/hooks/useGlobalLoading'

export function useProjects() {
  const teamStore = useTeamStore()
  const groupStore = useProjectGroupStore()
  const { selectedGroupRef } = storeToRefs(groupStore)
  const isLoading = ref(false)
  const projects = ref<ProjectAPI.ResponseProject[]>([])
  const selectProjectGroupRef = ref<InstanceType<typeof SelectProjectGroup>>()
  const router = useRouter()
  const { showGlobalLoading, hideGlobalLoading } = useGlobalLoading()
  // 跳转到项目详情
  const navigateToProjectDetail = async (project: ProjectAPI.ResponseProject) => {
    showGlobalLoading()

    await router.push({
      name: PROJECT_DETAIL_PATH_NAME,
      params: {
        project_id: project.id,
      },
    })

    hideGlobalLoading()
  }

  // 调整项目分组
  const showProjectGroupModal = (projectInfo: ProjectAPI.ResponseProject) => {
    selectProjectGroupRef.value?.show(projectInfo)
  }

  // 处理是否关注项目
  const handleFollowProject = async (project: ProjectAPI.ResponseProject) => {
    try {
      if (project.selfMember.isFollowed)
        await apiUnfollowProject(project.id)
      else await apiFollowProject(project.id as string)

      project.selfMember.isFollowed = !project.selfMember.isFollowed
      if (selectedGroupRef.value === 'followed')
        projects.value = await apiGetMyFollowedProjectList(teamStore.currentID)
    }
    catch (error) {
      //
    }
  }

  const loadPrjectListByGroupId = async () => {
    try {
      const groupKey = unref(selectedGroupRef)
      isLoading.value = true
      switch (groupKey) {
        case 'create':
          break
        case 'all':
          projects.value = await apiGetProjectList(teamStore.currentID)
          break
        case 'followed':
          projects.value = await apiGetMyFollowedProjectList(teamStore.currentID)
          break
        case 'my':
          projects.value = await apiGetMyProjectList(teamStore.currentID)
          break
        default:
          projects.value = await apiGetProjectListByGroupId(teamStore.currentID, groupKey as number, false)
          break
      }
    }
    catch (error) {
      projects.value = []
    }
    finally {
      isLoading.value = false
    }
  }

  return {
    isLoading,
    projects,
    selectProjectGroupRef,
    navigateToProjectDetail,
    showProjectGroupModal,
    refreshProjectList: loadPrjectListByGroupId,
    loadPrjectListByGroupId,
    handleFollowProject,
  }
}
