import useApi from '@/hooks/useApi'
import { ProjectInfo } from '@/typings'
import { getProjectDetailPath } from '@/router'
import { toggleFollowProject } from '@/api/project'
import { MemberAuthorityInProject } from '@/typings/member'
import { useProjectList } from '@/views/composables/useProjectList'

export const useProjects = (searchParams?: Record<string, any>) => {
  const { isLoading, projects } = useProjectList(searchParams)

  const router = useRouter()

  // 跳转到项目详情
  const goProjectDetail = (project: ProjectInfo) => {
    router.push(getProjectDetailPath(project.id))
  }

  // 处理创建项目
  const handleCreateProject = () => {}

  // 处理是否关注项目
  const handleFollowProject = async (project: ProjectInfo) => {
    try {
      await toggleFollowProject(project)
      project.is_followed = !project.is_followed
    } catch (error) {
      //
    }
  }

  return {
    isLoading,
    projects,

    handleFollowProject,
    goProjectDetail,
  }
}
