import { defineStore } from 'pinia'
import { isAxiosError } from 'axios'
import { apiChangeCurrentTeam, apiGetCurrentTeam, apiGetTeams } from '@/api/team'
import { NoTeamError } from '@/api/error'
import { Role } from '@/commons/constant'

export const useTeamStore = defineStore({
  id: 'team',
  state() {
    return {
      inited: false,
      hasTeam: false,
      teams: <TeamAPI.Team[]>[],
      teamIDIndex: <Record<string, TeamAPI.Team>>{},
      currentID: <string>'',
      currentRole: <TeamAPI.Role | undefined>undefined,
    }
  },
  getters: {
    isOwner: (state) => state.currentRole === Role.Owner,
    isMember: (state) => state.currentRole === Role.Member,
    isAdmin: (state) => state.currentRole === Role.Admin,
    currentTeam: (state) => state.teamIDIndex[state.currentID],
  },

  actions: {
    async init() {
      this.$state.inited = false
      const err = new NoTeamError()
      try {
        let teams: TeamAPI.Team[] = []
        let currentID: string = ''
        let currentRole: TeamAPI.Role | undefined
        const teamIDIndex: Record<string, TeamAPI.Team> = {}
        await Promise.all([
          (async () => {
            try {
              const res = await apiGetCurrentTeam()
              currentID = res.id
              currentRole = res.role
            } catch (e) {
              //
            }
          })(),

          apiGetTeams().then((res) => {
            for (const key in res.items) {
              const val = res.items[key]
              teamIDIndex[val.id!] = val
            }
            teams = res.items
          }),
        ])
        if (!teams.length) throw err

        this.$patch({
          teams,
          currentID,
          currentRole,
          teamIDIndex,
        })
        this.$state.inited = true
      } catch (e) {
        if (isAxiosError(e)) throw err
        throw e
      }
    },

    async switchTeam(teamID: string) {
      if (this.currentTeam?.id === teamID) return
      try {
        await apiChangeCurrentTeam(teamID)
        // this.currentID = teamID
        location.reload()
      } catch (error) {
        //
      }
    },
    async updateTeam(teamID: string, team: Partial<TeamAPI.Team>) {
      const target = this.teamIDIndex[teamID]
      for (const key in team) (<any>target)[<keyof typeof team>key] = team[<keyof typeof team>key]
    },
  },
})
