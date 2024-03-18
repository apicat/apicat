import type { Router } from 'vue-router'

import { setupAuthFilter } from './auth.filter'
import { setupGetUserInfoFilter } from './userInfo.filter'
import { setupGetTeamFilter } from './team.filter'
import { setupTitleI18n } from './titleI18n.filter'
import { setupNavigationFailureFilter } from './navigationFailures.filter'

export function setupRouterFilter(router: Router) {
  setupNavigationFailureFilter(router)
  setupTitleI18n(router)
  setupAuthFilter(router)
  setupGetUserInfoFilter(router)
  setupGetTeamFilter(router)
}
