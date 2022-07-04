import { getProjectDetail, getProjectStatus, settingProject } from '@/api/project'
import { PROJECT_ROLES_KEYS, PROJECT_VISIBLE_TYPES, PROJECT_VISIBLE_LIST } from '@/common/constant'
import { defineStore } from 'pinia'
import { Storage } from '@natosoft/shared'

interface ProjectState {
    projectInfo: any
    projectAuthInfo: any
}

export const useProjectStore = defineStore({
    id: 'project',

    state: (): ProjectState => ({
        projectInfo: null,
        projectAuthInfo: {},
    }),

    getters: {
        isManager: (state) => state.projectAuthInfo && state.projectAuthInfo.authority === PROJECT_ROLES_KEYS.MANAGER,
        isDeveloper: (state) => state.projectAuthInfo && state.projectAuthInfo.authority === PROJECT_ROLES_KEYS.DEVELOPER,
        isReader: (state) => state.projectAuthInfo && state.projectAuthInfo.authority === PROJECT_ROLES_KEYS.READER,
        isGuest: (state) => state.projectAuthInfo && state.projectAuthInfo.authority === PROJECT_ROLES_KEYS.NONE,
        isPrivate: (state) => state.projectAuthInfo && state.projectAuthInfo.visibility === PROJECT_VISIBLE_TYPES.PRIVATE,
    },

    actions: {
        async getProjectDetail(project_id: number) {
            const token = Storage.get(Storage.KEYS.SECRET_PROJECT_TOKEN + project_id || '', true)
            const { data } = await getProjectDetail(project_id, token)
            // todo temp code
            this.projectInfo = { ...(data || {}), visible: PROJECT_VISIBLE_LIST.find((item) => item.key === data.visibility)?.value }
            return this.projectInfo
        },

        async getProjectAuth(project_id: number) {
            const { data } = await getProjectStatus(project_id)
            this.projectAuthInfo = data || {}
            this.projectAuthInfo.id = project_id
            // 是否在项目中
            this.projectAuthInfo.in_this = this.projectAuthInfo.authority !== PROJECT_ROLES_KEYS.NONE
            // console.log('获取项目权限信息：', JSON.stringify(this.projectAuthInfo))
            return this.projectAuthInfo
        },

        async updateProjectInfo(project: any) {
            const res = await settingProject(project)
            // todo temp code
            this.projectInfo = { ...this.projectInfo, ...project, visibility: PROJECT_VISIBLE_LIST.find((item) => item.value === project.visible)?.key }
            return res
        },

        clearProjectInfo() {
            this.projectInfo = null
        },
    },
})
