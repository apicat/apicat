import { getProjectList, getProjectDetail, getProjectServerUrlList, saveProjectServerUrlList } from '@/api/project'
import { ProjectListCoverBgColors, ProjectListCoverIcons } from '@/commons'
import { ProjectInfo } from '@/typings/project'
import { getProjectDefaultCover } from '@/views/project/logic/useProjectCover'
import { defineStore } from 'pinia'

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
  },
  actions: {
    async getProjects() {
      const projects: any = await getProjectList()
      this.projects = projects
    },

    async getProjectDetailInfo(project_id: string): Promise<ProjectInfo> {
      const project = await getProjectDetail(project_id)
      return project as any
    },

    setCurrentProjectInfo(info: ProjectInfo) {
      this.projectDetailInfo = { ...this.projectDetailInfo, ...info }
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
