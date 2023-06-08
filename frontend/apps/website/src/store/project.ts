import { getProjectList, getProjectDetail, getProjectServerUrlList, saveProjectServerUrlList } from '@/api/project'
import { ProjectListCoverBgColors, ProjectListCoverIcons } from '@/commons'
import { MemberAuthorityInProject, MemberAuthorityMap } from '@/typings/member'
import { ProjectInfo } from '@/typings/project'
import { getProjectDefaultCover } from '@/views/project/logic/useProjectCover'
import { defineStore } from 'pinia'
import { pinia } from '@/plugins'

interface ProjectState {
  projects: ProjectInfo[]
  projectDetailInfo: ProjectInfo | null
  urlServers: Array<any>
}

export const uesProjectStore = defineStore('project', {
  state: (): ProjectState => ({
    projects: [],
    projectDetailInfo: null,
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
  },
})

export default uesProjectStore

export const uesProjectStoreWithOut = () => uesProjectStore(pinia)
