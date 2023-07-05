import type { Router } from 'vue-router'
import { setupGetProjectInfoFilter } from './projectInfo.filter'
import { setupMemberPermissionFilter } from './memberPermission.filter'
import { setupAuthFilter } from './auth.filter'
import { setupGetUserInfoFilter } from './userInfo.filter'
import { setupShareFilter } from './share.filter'

export const setupRouterFilter = (router: Router) => {
  setupAuthFilter(router)
  setupGetUserInfoFilter(router)
  setupGetProjectInfoFilter(router)
  setupMemberPermissionFilter(router)
  setupShareFilter(router)
}
