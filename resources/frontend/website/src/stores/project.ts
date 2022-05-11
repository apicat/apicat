import { getProjectDetail, settingProject } from '@/api/project'
import { PROJECT_ROLES_MAP } from '@ac/shared'
import { defineStore } from 'pinia'
interface ProjectState {
    projectInfo: any
}

export const useProjectStore = defineStore({
    id: 'project',

    state: (): ProjectState => ({
        projectInfo: null,
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
