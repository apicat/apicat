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

  actions: {
    async getProjectGroups() {
      try {
        const groups = await getProjectGroupList()
        this.projectGroups = groups
      } catch (error) {
        this.projectGroups = []
      }

      this.projectGroups = [
        { id: 1, name: '1' },
        { id: 2, name: '2' },
        { id: 3, name: '3' },
      ]
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
