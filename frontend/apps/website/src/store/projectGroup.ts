import { defineStore, storeToRefs } from 'pinia'
import { useI18n } from 'vue-i18n'
import { apiSortProjectGroup, createProjectGroup, deleteProjectGroup, getProjectGroupList, updateProjectGroup } from '@/api/project/group'
import { useTeamStore } from '@/store/team'
import { Storage } from '@/commons'

interface ProjectGroupState {
  t: any
  projectGroups: ProjectAPI.ResponseGroup[]
  selectedGroupRef: ProjectGroupSelectKey
}

export const useProjectGroupStore = defineStore('projectGroup', {
  state: (): ProjectGroupState => {
    const { t } = useI18n()
    return {
      t,
      projectGroups: [],
      selectedGroupRef: Storage.get(Storage.KEYS.SELECTED_PROJECT_GROUP) || 'all',
    }
  },

  getters: {
    groupsForOptions: (state): ProjectAPI.ResponseGroup[] => [
      { name: state.t('app.project.groups.noGroup'), id: 0 },
      ...state.projectGroups,
    ],
  },

  actions: {
    saveGroupKeyToStorage(key: ProjectGroupSelectKey) {
      if (key !== 'create')
        Storage.set(Storage.KEYS.SELECTED_PROJECT_GROUP, key)
    },

    async getProjectGroups(): Promise<ProjectAPI.ResponseGroup[]> {
      try {
        const { currentID } = storeToRefs(useTeamStore())
        const groups = await getProjectGroupList(currentID.value)
        this.projectGroups = groups
      }
      catch (error) {
        this.projectGroups = []
      }

      return this.projectGroups
    },

    async createOrUpdateProjectGroup(group: ProjectAPI.ResponseGroup) {
      const teamStore = useTeamStore()
      const res = group.id
        ? await updateProjectGroup(group.id as string, group)
        : await createProjectGroup(teamStore.currentID, {
          name: group.name,
          teamID: teamStore.currentID,
        })

      if (group.id) {
        this.projectGroups = this.projectGroups.map((val) => {
          if (val.id === group.id)
            return group
          return val
        })
      }
      else {
        res && this.projectGroups.push(res)
      }

      return res
    },

    async deleteProjectGroup(group: ProjectAPI.ResponseGroup) {
      await deleteProjectGroup(group.id! as number)
      this.projectGroups = this.projectGroups.filter((val: any) => val.id !== group.id)
    },

    async sortGroup() {
      const teamStore = useTeamStore()
      await apiSortProjectGroup(teamStore.currentTeam.id, {
        groupIDs: this.projectGroups.map((val: any) => val.id),
      })
    },
  },
})

export default useProjectGroupStore
