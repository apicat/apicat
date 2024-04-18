import { defineStore } from 'pinia'
import { apiChangeCurrentTeam, apiGetCurrentTeam, apiGetTeams } from '@/api/team'
import { Role } from '@/commons/constant'

export const useTeamStore = defineStore({
  id: 'team',

  state() {
    return {
      teams: <TeamAPI.Team[]>[],
      teamMaps: new Map<string, TeamAPI.Team>(),
      currentID: <string>'',
      currentRole: <TeamAPI.Role | undefined>undefined,
    }
  },

  getters: {
    isOwner: state => state.currentRole === Role.Owner,
    isMember: state => state.currentRole === Role.Member,
    isAdmin: state => state.currentRole === Role.Admin,
    currentTeam: state => state.teamMaps.get(state.currentID),
    hasTeam: state => state.teamMaps.size > 0,
  },

  actions: {
    async init() {
      await this.getTeams()
      this.hasTeam && await this.getActiveTeam()
    },

    async getActiveTeam() {
      const { id, role } = await apiGetCurrentTeam()
      this.currentID = id
      this.currentRole = role
    },

    async getTeams() {
      const { items } = await apiGetTeams()
      this.teams = items || []
      this.teamMaps = new Map(items.map(team => [team.id!, team]))
    },

    async switchTeam(teamID: string) {
      if (this.currentTeam?.id === teamID)
        return
      try {
        await apiChangeCurrentTeam(teamID)
        location.reload()
      }
      catch (error) {
        //
      }
    },

    async updateTeam(teamID: string, team: Partial<TeamAPI.Team>) {
      const target = this.teamMaps.get(teamID)
      for (const key in team) (<any>target)[<keyof typeof team>key] = team[<keyof typeof team>key]
    },
  },
})
