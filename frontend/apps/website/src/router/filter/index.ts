import type { Router } from 'vue-router'
import { setupGetProjectInfoFilter } from './projectInfo.filter'
import { setupMemberPermissionFilter } from './memberPermission.filter'
import { setupAuthFilter } from './auth.filter'

export const setupRouterFilter = (router: Router) => {
  setupAuthFilter(router)
  setupGetProjectInfoFilter(router)
  setupMemberPermissionFilter(router)
}
