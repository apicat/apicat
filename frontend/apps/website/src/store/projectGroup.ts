import { getProjectGroupList, createProjectGroup, updateProjectGroup, deleteProjectGroup } from '@/api/projectGroup'
import { ProjectGroup } from '@/typings'
import { defineStore } from 'pinia'

interface ProjectGroupState {
  projectGroups: ProjectGroup[]
}

export const useProjectGroupStore = defineStore('projectGroup', {
  state: (): ProjectGroupState => ({
    projectGroups: [],
  }),

  getters: {
    groupsForOptions: (state): ProjectGroup[] => [{ name: '不分组', id: 0 }, ...state.projectGroups],
  },

  actions: {
    async getProjectGroups() {
      try {
        const groups = await getProjectGroupList()
        this.projectGroups = groups
      } catch (error) {
        this.projectGroups = []
      }
      return this.projectGroups
    },

    async createOrUpdateProjectGroup(group: ProjectGroup) {
      group.id ? await updateProjectGroup(group) : await createProjectGroup(group)
    },

    async deleteProjectGroup(id: number) {
      await deleteProjectGroup(id)
    },
  },
})

export default useProjectGroupStore
