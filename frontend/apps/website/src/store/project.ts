import { getProjectList, getProjectDetail, getProjectServerUrlList, saveProjectServerUrlList, getProjectAuthInfo } from '@/api/project'
import { ProjectListCoverBgColors, ProjectListCoverIcons, ProjectVisibilityEnum } from '@/commons'
import { MemberAuthorityInProject, MemberAuthorityMap } from '@/typings/member'
import { ProjectInfo } from '@/typings/project'
import { getProjectDefaultCover } from '@/views/project/logic/useProjectCover'
import { defineStore } from 'pinia'
import { pinia } from '@/plugins'

interface ProjectAuthInfo {
  inThisProject: boolean
  hasShared: boolean
  isPrivate: boolean
}
interface ProjectState {
  projects: ProjectInfo[]
  projectDetailInfo: ProjectInfo | null
  projectAuthInfo: ProjectAuthInfo | null
  urlServers: Array<any>
}

export const useProjectStore = defineStore('project', {
  state: (): ProjectState => ({
    projects: [],
    projectDetailInfo: null,
    projectAuthInfo: null,
    urlServers: [],
  }),

  getters: {
    projectList: (state) =>
      state.projects.map((info) => {
        try {
          info.cover = JSON.parse(info.cover as string)
        } catch (error) {
          info.cover = getProjectDefaultCover({ coverBgColor: ProjectListCoverBgColors[1], coverIcon: ProjectListCoverIcons[0], type: 'icon' })
        }
        return info
      }),

    projectAuths: () => {
      return Object.keys(MemberAuthorityMap)
        .filter((key: string) => key !== MemberAuthorityInProject.MANAGER)
        .map((key: string) => {
          return {
            text: (MemberAuthorityMap as any)[key],
            value: key,
          }
        })
    },

    isManager: (state) => state.projectDetailInfo?.authority === MemberAuthorityInProject.MANAGER,

    isWriter: (state) => state.projectDetailInfo?.authority === MemberAuthorityInProject.WRITE,

    isReader: (state) => state.projectDetailInfo?.authority === MemberAuthorityInProject.READ,

    isPrivate: (state) => state.projectDetailInfo?.visibility === ProjectVisibilityEnum.PRIVATE,
  },
  actions: {
    async getProjects() {
      const projects: any = await getProjectList()
      this.projects = projects
    },

    async getProjectDetailInfo(project_id: string): Promise<ProjectInfo> {
      const project = await getProjectDetail(project_id)
      this.setCurrentProjectInfo(project as any)
      return project as any
    },

    setCurrentProjectInfo(info?: ProjectInfo) {
      this.projectDetailInfo = info ? { ...this.projectDetailInfo, ...info } : null
    },

    clearCurrentProjectInfo() {
      this.setCurrentProjectInfo()
    },

    async getUrlServers(project_id: string) {
      const urls: any = await getProjectServerUrlList(project_id)
      this.urlServers = urls
      return urls
    },

    async saveProjectServerUrlListApi({ project_id, urls }: any) {
      await saveProjectServerUrlList({ project_id, urls })
      this.urlServers = urls
    },

    async getProjectAuthInfo(project_id: string): Promise<ProjectAuthInfo> {
      const { authority, visibility, secret_key } = await getProjectAuthInfo(project_id)
      this.projectAuthInfo = {
        inThisProject: authority !== MemberAuthorityInProject.NONE,
        hasShared: !!secret_key,
        isPrivate: visibility === ProjectVisibilityEnum.PRIVATE,
      }
      return this.projectAuthInfo
    },
  },
})

export default useProjectStore

export const useProjectStoreWithOut = () => useProjectStore(pinia)
