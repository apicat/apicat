import type { Router } from 'vue-router'
import { setupGetProjectInfoFilter, setupGetProjectAuthInfoFilter } from './projectInfo.filter'
import { setupMemberPermissionFilter } from './memberPermission.filter'
import { setupAuthFilter } from './auth.filter'
import { setupGetUserInfoFilter } from './userInfo.filter'
import { setupShareDocumentFilter } from './share.filter'
import { setupDocumentDetailFilter } from './document.detail.filter'

export const setupRouterFilter = (router: Router) => {
  setupAuthFilter(router)
  setupGetUserInfoFilter(router)

  setupGetProjectAuthInfoFilter(router)
  setupDocumentDetailFilter(router)

  setupGetProjectInfoFilter(router)

  setupMemberPermissionFilter(router)

  setupShareDocumentFilter(router)
}
