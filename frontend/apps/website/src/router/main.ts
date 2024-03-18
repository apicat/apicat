import type { RouteRecordRaw } from 'vue-router'
import { MAIN_PATH, MAIN_PATH_NAME } from './constant'
import { pojectsRoute } from './projects'
import { noTeamRoute, teamRoute } from './team'
import { iterationsRoute } from './iterations'
import { userRoute } from './userSetting'

import MainLayout from '@/layouts/MainLayout.vue'

export const mainRoute: RouteRecordRaw = {
  name: MAIN_PATH_NAME,
  path: MAIN_PATH,
  redirect: '/projects',
  component: MainLayout,
  children: [pojectsRoute, iterationsRoute, userRoute, teamRoute, noTeamRoute],
}
