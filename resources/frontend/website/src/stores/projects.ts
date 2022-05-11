import { getProjectGroupList, createProjectGroup, renameProjectGroup, removeProjectGroup, sortProjectGroup } from '@/api/project'
import { Storage } from '@ac/shared'
import { defineStore } from 'pinia'

interface ProjectState {
    projectGroupList: Array<any>
    activeGroup: any
}

const DEFAULT_ACTIVE_GROUP = { id: 0, name: '所有项目' }

export const useProjectsStore = defineStore({
    id: 'projects',

    state: (): ProjectState => ({
        projectGroupList: [],
        activeGroup: Storage.get(Storage.KEYS.ACTIVE_PROJECT_GROUP) || DEFAULT_ACTIVE_GROUP,
    }),

    getters: {},

    actions: {
        // 获取项目列表
        async getProjectGroupList() {
            const { data = [] } = (await getProjectGroupList()) || {}
            this.projectGroupList = data
        },

        // 切换项目分组
        switchProjectGroup(group = DEFAULT_ACTIVE_GROUP) {
            Storage.set(Storage.KEYS.ACTIVE_PROJECT_GROUP, group)
            this.activeGroup = group
        },

        async addProjectGroup(group: any) {
            const res = await createProjectGroup(group)
            res && this.projectGroupList.push(res.data)
            return res
        },

        async renameProjectGroup(newGroup: any) {
            const res = await renameProjectGroup(newGroup)
            if (res) {
                const oldGroup = this.projectGroupList.find((item) => item.id == newGroup.id)
                oldGroup.name = newGroup.name
            }
            return res
        },

        async deleteProjectGroup(group: any) {
            const res = await removeProjectGroup(group.id)
            this.projectGroupList = this.projectGroupList.filter((item) => item.id !== group.id)
            // 删除当前已选中分组
            if (this.activeGroup.id === group.id) {
                this.switchProjectGroup()
            }
            return res
        },

        async sortProjectGroup(oldItemIndex: number, newItemIndex: number) {
            const changeItem = this.projectGroupList.splice(oldItemIndex, 1)[0]
            this.projectGroupList.splice(newItemIndex, 0, changeItem)

            await sortProjectGroup(this.projectGroupList.map((item) => item.id))
        },
    },
})
