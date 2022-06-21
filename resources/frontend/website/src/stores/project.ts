import { getProjectDetail, settingProject } from '@/api/project'
import { PROJECT_ROLES_MAP } from '@/common/constant'
import { defineStore } from 'pinia'
import delay from 'delay'

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
        isManager: (state) => state.projectInfo && state.projectInfo.authority === PROJECT_ROLES_MAP.MANAGER,
        isDeveloper: (state) => state.projectInfo && state.projectInfo.authority === PROJECT_ROLES_MAP.DEVELOPER,
        isReader: (state) => state.projectInfo && state.projectInfo.authority === PROJECT_ROLES_MAP.READER,
    },

    actions: {
        async getProjectDetail(project_id: number) {
            const { data } = await getProjectDetail(project_id)
            this.projectInfo = data || null
            return this.projectInfo
        },

        async getProjectAuth(project_id: number) {
            // const { data } = await getProjectDetail(project_id)
            // this.projectInfo = data || null
            console.log('获取项目权限信息：', project_id)
            await delay(2000)
            this.projectAuthInfo = {
                id: project_id,
            }
            return this.projectAuthInfo
        },

        async updateProjectInfo(project: any) {
            const res = await settingProject(project)
            this.projectInfo = { ...this.projectInfo, ...project }
            return res
        },

        clearProjectInfo() {
            this.projectInfo = null
        },
    },
})
