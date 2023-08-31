import { ProjectGroupSelectKey, ProjectInfo } from '@/typings'
import SelectProjectGroup from '../components/SelectProjectGroup.vue'
import { getProjectDetailPath } from '@/router'
import { getMyFollowedProjectList, getMyProjectList, getProjectList, getProjectListByGroupId, toggleFollowProject } from '@/api/project'

export const useProjects = (selectedGroupRef: Ref<ProjectGroupSelectKey>) => {
  const router = useRouter()
  const isLoading = ref(false)
  const projects = ref<ProjectInfo[]>([])
  const selectProjectGroupRef = ref<InstanceType<typeof SelectProjectGroup>>()

  // 跳转到项目详情
  const goProjectDetail = (project: ProjectInfo) => {
    router.push(getProjectDetailPath(project.id))
  }

  // 调整项目分组
  const changeProjectGroup = (projectInfo: ProjectInfo) => {
    selectProjectGroupRef.value?.show(projectInfo)
  }

  // 处理是否关注项目
  const handleFollowProject = async (project: ProjectInfo) => {
    try {
      await toggleFollowProject(project)
      project.is_followed = !project.is_followed
    } catch (error) {
      //
    }
  }

  const loadPrjectListByGroupId = async () => {
    try {
      const groupKey = unref(selectedGroupRef)
      isLoading.value = true
      switch (groupKey) {
        case 'all':
          projects.value = await getProjectList()
          break
        case 'followed':
          projects.value = await getMyFollowedProjectList()
          break
        case 'my':
          projects.value = await getMyProjectList()
          break
        default:
          projects.value = await getProjectListByGroupId(groupKey as number)
          break
      }
    } catch (error) {
      projects.value = []
    } finally {
      isLoading.value = false
    }
  }

  // 检测项目分组变化
  watch(selectedGroupRef, async () => await loadPrjectListByGroupId(), { immediate: true })

  return {
    isLoading,
    projects,
    selectProjectGroupRef,
    refreshProjectList: loadPrjectListByGroupId,
    handleFollowProject,
    goProjectDetail,
    changeProjectGroup,
  }
}
