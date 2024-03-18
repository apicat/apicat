import DefaultAjax from '../Ajax'

/**
 * get all team current user is in
 */
export async function apiGetTeams(): Promise<TeamAPI.ResponseTeams> {
  return DefaultAjax.get('/teams')
}

/**
 * get current user's current team to display
 */
export async function apiGetCurrentTeam(): Promise<TeamAPI.ResponseCurrentTeam> {
  return DefaultAjax.get('/teams/current', undefined, { isShowErrorMsg: false })
}

/**
 * create team
 */
export async function apiCreateTeam(data: TeamAPI.RequestCreateTeam): Promise<TeamAPI.ResponseCreateTeam> {
  return DefaultAjax.post('/teams', data)
}

/**
 * change current user's current team to display
 */
export async function apiChangeCurrentTeam(teamID: string): Promise<void> {
  return DefaultAjax.put(`/teams/${teamID}/switch`)
}

/**
 * get team members by `teamID`
 */
export async function apiGetTeamMembers(
  params: Partial<TeamAPI.RequestMembers>,
  teamID: string,
): Promise<GlobalAPI.ResponseTable<TeamAPI.TeamMember[]>> {
  return DefaultAjax.get(`/teams/${teamID}/members`, {
    params,
  })
}

/**
 * 获取团队成员列表（只获取`status`为`active`的成员）
 */
export async function apiGetTeamActiveMembers(
  params: Partial<TeamAPI.RequestMembers>,
  teamID: string,
): Promise<GlobalAPI.ResponseTable<TeamAPI.TeamMember[]>> {
  params = { ...(params || {}), status: 'active' } as any
  return DefaultAjax.get(`/teams/${teamID}/members`, {
    params,
  })
}

/**
 * edit one member's info. (only `role` supported for now)
 */
export async function apiEditMember(
  teamID: string,
  memberID: number,
  data: TeamAPI.RequestEditMember,
): Promise<TeamAPI.Team> {
  return DefaultAjax.put(`/teams/${teamID}/members/${memberID}`, data)
}

// 删除团队成员
export function apiRemoveTeamMember(teamID: string, memberID: number) {
  return DefaultAjax.delete(`/teams/${teamID}/members/${memberID}`, null, {
    isShowErrorMsg: true,
    isShowSuccessMsg: true,
  })
}

// 退出团队
export function apiQuitTeam(teamID: string) {
  return DefaultAjax.delete(`/teams/${teamID}/members`, null, {
    isShowErrorMsg: true,
    isShowSuccessMsg: true,
  })
}

/**
 * get invite token (for team invitation)
 */
export async function apiGetInviteToken(teamID: string): Promise<TeamAPI.ResponseInvite> {
  return DefaultAjax.get(`/teams/${teamID}/invitation-tokens`)
}

/**
 * reset invite token (won't return new token)
 */
export async function apiResetInviteToken(teamID: string): Promise<TeamAPI.ResponseInvite> {
  return DefaultAjax.put(`/teams/${teamID}/invitation-tokens`)
}

/**
 * (`owner`)update team setting (currently only support name)
 */
export async function apiTeamSetting(teamID: string, data: TeamAPI.RequestTeamSetting): Promise<void> {
  return DefaultAjax.put(`/teams/${teamID}/setting`, data, { isShowSuccessMsg: true })
}

/**
 * (`owner`)transfer ownership to an admin
 */
export async function apiTransferOwnership(teamID: string, data: TeamAPI.RequestTransferOwnership): Promise<void> {
  return DefaultAjax.put(`/teams/${teamID}/transfer`, data)
}

/**
 * (`owner`)delete team
 */
export async function apiDeleteTeam(teamID: string): Promise<void> {
  return DefaultAjax.delete(`/teams/${teamID}`)
}

/**
 * join a team (given invitation token)
 */
export async function apiJoinTeam(invitationToken: string): Promise<void> {
  return DefaultAjax.post('/teams/join', { invitationToken })
}

// =================================================
// join team get team detail
export interface InviteTokenTeamData {
  team: string
  inviter: string
}

export function getJoinTeamInfoByToken(invitationToken: string): Promise<InviteTokenTeamData> {
  return DefaultAjax.get(`/teams/invitation-tokens/${invitationToken}`, {}, { isShowErrorMsg: false })
}
