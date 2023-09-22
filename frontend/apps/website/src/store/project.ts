import { getProjectDetail, getProjectServerUrlList, saveProjectServerUrlList, updateProjectBaseInfo } from '@/api/project'
import { Cookies, ProjectVisibilityEnum } from '@/commons'
import { MemberAuthorityInProject, MemberAuthorityMap } from '@/typings/member'
import { ProjectInfo } from '@/typings/project'
import { defineStore } from 'pinia'
import { pinia } from '@/plugins'
import { getProjectAuthInfo } from '@/api/shareProject'
import { getProjectDetailPath } from '@/router/project.detail'
import { PROJECT_DETAIL_PATH_NAME } from '@/router/constant'

interface ProjectAuthInfo {
  project_id: string
  inThisProject: boolean
  hasShared: boolean
  isPrivate: boolean
}
interface ProjectState {
  projectDetailInfo: ProjectInfo | null
  projectAuthInfo: ProjectAuthInfo | null
  urlServers: Array<any>
  isShowProjectSecretLayer: boolean
}

export const useProjectStore = defineStore('project', {
  state: (): ProjectState => ({
    projectDetailInfo: null,
    projectAuthInfo: null,
    urlServers: [],
    isShowProjectSecretLayer: false,
  }),

  getters: {
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
    isGuest: (state) => state.projectDetailInfo?.authority === MemberAuthorityInProject.NONE,
    isPrivate: (state) => state.projectDetailInfo?.visibility === ProjectVisibilityEnum.PRIVATE,

    hasInputSecretKey: (state) => !!(Cookies.get(Cookies.KEYS.SHARE_PROJECT + state.projectAuthInfo?.project_id) || ''),

    isProjectRoute: function (): boolean {
      return !!this.$router.currentRoute.value.matched.find((item) => item.name === PROJECT_DETAIL_PATH_NAME)
    },
  },
  actions: {
    async getProjectDetailInfo(project_id: string): Promise<ProjectInfo> {
      const token = Cookies.get(Cookies.KEYS.SHARE_PROJECT + project_id)
      const project = await getProjectDetail(project_id, token ? { token } : {})
      this.updateCurrentProjectInfo(project)
      return project
    },

    updateCurrentProjectInfo(info?: ProjectInfo) {
      this.projectDetailInfo = info ? { ...this.projectDetailInfo, ...info } : null
    },

    clearCurrentProjectInfo() {
      this.updateCurrentProjectInfo()
    },

    async updateProectInfo(info: ProjectInfo) {
      await updateProjectBaseInfo(info)
      this.updateCurrentProjectInfo(info)
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
      const { authority, visibility, has_shared } = await getProjectAuthInfo(project_id)
      this.projectAuthInfo = {
        project_id,
        inThisProject: authority !== MemberAuthorityInProject.NONE,
        hasShared: has_shared,
        isPrivate: visibility === ProjectVisibilityEnum.PRIVATE,
      }
      return this.projectAuthInfo
    },
    switchProjectSecretLayer(isShow = true) {
      this.isShowProjectSecretLayer = isShow
    },
    showProjectSecretLayer() {
      this.switchProjectSecretLayer()
    },
    hideProjectSecretLayer() {
      this.switchProjectSecretLayer(false)
    },
    removeProjectSecretKeyWithReload() {
      const { project_id } = this.projectAuthInfo!
      Cookies.remove(Cookies.KEYS.SHARE_PROJECT + project_id)
      setTimeout(() => location.replace(getProjectDetailPath(project_id as string)), 500)
    },
  },
})

export default useProjectStore

export const useProjectStoreWithOut = () => useProjectStore(pinia)
