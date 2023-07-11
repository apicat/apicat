import type { Router } from 'vue-router'
import { setupGetProjectInfoFilter, setupGetProjectAuthInfoFilter } from './projectInfo.filter'
import { setupMemberPermissionFilter } from './memberPermission.filter'
import { setupAuthFilter } from './auth.filter'
import { setupGetUserInfoFilter } from './userInfo.filter'
import { setupShareDocumentFilter, setupShareProjectDetailFilter } from './share.filter'

export const setupRouterFilter = (router: Router) => {
  setupAuthFilter(router)
  setupGetUserInfoFilter(router)

  setupGetProjectAuthInfoFilter(router)

  setupShareProjectDetailFilter(router)

  setupGetProjectInfoFilter(router)

  setupMemberPermissionFilter(router)

  setupShareDocumentFilter(router)
}
