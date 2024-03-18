declare namespace TeamAPI {
  type Status = import('@/commons/constant').Status
  type Role = import('@/commons/constant').Role
  type EditableRole = import('@/commons/constant').EditableRole

  interface Team {
    id: string
    name: string
    createdAt: string
    updatedAt: string
    avatar?: string
    membersCount?: number
    owner?: UserAPI.ResponseUserInfo
    role?: Role
  }
  interface TeamItem extends TeamAPI.Team {
    text: string
    onClick: any
  }
  interface ResponseTeams {
    items: Team[]
  }
  type ResponseCurrentTeam = Team

  interface RequestCreateTeam {
    name: string
    avatar?: string
  }
  interface ResponseCreateTeam {
    name: string
    avatar?: string
    createdAt?: string
    membersCount?: number
    owner?: UserAPI.ResponseUserInfo
    teamID?: string
  }

  interface RequestMembers extends GlobalAPI.RequestTable {
    roles: Role
  }

  interface TeamMember {
    id: number
    role: Role
    teamID: string
    createdAt: string
    updatedAt: string
    user: UserData
    status?: Status
  }
  interface UserData {
    name: string
    email: string
    avatar?: string
    level?: number
  }

  interface RequestEditMember {
    role?: 'admin' | 'member'
    status?: Status
  }
  type ResponseEditMember = TeamMember

  interface ResponseInvite {
    invitationToken: string
  }

  interface RequestTeamSetting {
    name: string
    avatar?: string
  }
  interface RequestTransferOwnership {
    memberID: number
  }
}
