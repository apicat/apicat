import { defineStore } from 'pinia'
import useLocaleStore from './locale'
import { pinia } from '@/plugins'
import { PROJECT_DETAIL_PATH_NAME } from '@/router/constant'
import { Authority, MemberAuthorityMap, Visibility } from '@/commons/constant'
import { apiDeleteProject, apiGetProject, apiSetProjectGeneral } from '@/api/project'
import { apiTransferProject } from '@/api/project/index'
import { setProjectSharedToken } from '@/api/shareToken'
import { apiGetProjectShareStatus, apiSendProjectSharekey } from '@/api/project/share'
import { BadRequestError } from '@/api/error'

interface ProjectState {
  // i18n
  t: any
  project: ProjectAPI.ResponseProject | undefined
  //  项目权限信息（项目分享状态）
  projectAuthInfo: ShareAPI.ProjectAuthInfo | undefined
  /* 分享项目密钥层 - 字符串为空则代表没有 */
  isShowProjectSecretLayer: boolean
}

export const useProjectStore = defineStore('project', {
  state: (): ProjectState => {
    const { t } = useLocaleStore()
    return {
      t,
      projectAuthInfo: undefined,
      isShowProjectSecretLayer: false,
      project: undefined,
    }
  },

  getters: {
    projectAuths: (state) => {
      return Object.keys(MemberAuthorityMap)
        .filter((key: string) => key !== Authority.Manage)
        .map((key: string) => {
          return {
            text: state.t((MemberAuthorityMap as any)[key]),
            value: key,
          }
        })
    },

    isManager: state => state.project?.selfMember.permission === Authority.Manage,
    isWriter: state => state.project?.selfMember.permission === Authority.Write,
    isReader: state => state.project?.selfMember.permission === Authority.Read,
    isGuest: state => state.project?.selfMember.permission === Authority.None,
    isPrivate: state => state.project?.visibility === Visibility.Private,

    isProjectRoute(): boolean {
      return !!this.$router.currentRoute.value.matched.find(item => item.name === PROJECT_DETAIL_PATH_NAME)
    },
    projectID: state => state.project?.id,
    mockURL: state => state.project?.mockURL,
  },
  actions: {
    async getProjejctAuthInfo(projectID: string) {
      try {
        this.projectAuthInfo = await apiGetProjectShareStatus(projectID)
        this.projectAuthInfo.projectID = projectID
      }
      catch (error) {
        //
      }
      return this.projectAuthInfo
    },

    async getProjectInfoById(projectID: string) {
      this.project = await apiGetProject(projectID)
    },

    async clearProject() {
      this.projectAuthInfo = undefined
      this.project = undefined
    },

    async updateProjectGeneral({ id, title, visibility, cover, description }: ProjectAPI.ResponseProject) {
      const data: ProjectAPI.RequestSetProjectGeneral = {
        title,
        visibility,
        cover: cover as string,
        description,
      }
      const res = await apiSetProjectGeneral(id, data)
      this.$patch({ project: { ...this.project, ...data } })
      return res
    },

    refreshProject(projectID: string) {
      return this.getProjectInfoById(projectID)
    },

    async deleteProject(projectID: string) {
      const a = await apiDeleteProject(projectID)
      if (projectID === this.project?.id)
        this.clearProject()
      return a
    },

    async deleteCurrentProject() {
      const a = await apiDeleteProject(this.project!.id)
      this.clearProject()
      return a
    },

    async transferProject(memberID: number) {
      await apiTransferProject(this.project!.id, memberID)
      await this.getProjectInfoById(this.project!.id)
    },

    // 通过密钥获取项目分享token
    async getProjectShareTokenBySecretCode(code: string) {
      if (!this.projectAuthInfo)
        throw new Error(this.t('app.project.share.noAuthInfo'))

      try {
        const projectID = this.projectAuthInfo.projectID!
        const res = await apiSendProjectSharekey(projectID, code)
        setProjectSharedToken(projectID, res.shareCode)
        await this.getProjectInfoById(projectID)
        this.isShowProjectSecretLayer = false
      }
      catch (e) {
        if (e instanceof BadRequestError)
          this.isShowProjectSecretLayer = true
      }
    },
  },
})

export default useProjectStore

export const useProjectStoreWithOut = () => useProjectStore(pinia)
