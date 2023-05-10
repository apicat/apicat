import type { Router } from 'vue-router'
import { setupGetProjectInfoFilter } from './getProjectInfo'
import { setupAuthFilter } from './auth.filter'

export const setupRouterFilter = (router: Router) => {
  setupAuthFilter(router)
  setupGetProjectInfoFilter(router)
}
