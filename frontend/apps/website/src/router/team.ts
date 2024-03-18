import { compile } from 'path-to-regexp'
import {
  TEAM_CREATE_NAME,
  TEAM_CREATE_PATH,
  TEAM_JOIN_NAME,
  TEAM_JOIN_PATH,
  TEAM_NAME,
  TEAM_PATH,
} from './constant'

export function getTeamJoinPath(token: string) {
  return location.origin + compile(TEAM_JOIN_PATH)({ token })
}

export const noTeamRoute = {
  path: TEAM_CREATE_PATH,
  name: TEAM_CREATE_NAME,
  meta: { title: 'app.pageTitles.createTeam' },
  component: () => import('@/views/team/noTeam.vue'),
}

export const teamRoute = {
  path: TEAM_PATH,
  name: TEAM_NAME,
  meta: { title: 'app.pageTitles.team' },
  component: () => import('@/views/team/Team.vue'),
}

export const joinTeamRoute = {
  path: TEAM_JOIN_PATH,
  name: TEAM_JOIN_NAME,
  meta: { ignoreAuth: true, title: 'app.pageTitles.joinTeam' },
  component: () => import('@/views/team/JoinTeam.vue'),
}
