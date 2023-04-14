import type { Router } from 'vue-router'
import { setupGetProjectInfoFilter } from './getProjectInfo'

export const setupRouterFilter = (router: Router) => {
  setupGetProjectInfoFilter(router)
}
