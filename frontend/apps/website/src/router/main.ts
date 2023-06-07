import { RouteRecordRaw } from 'vue-router'
import MainLayout from '@/layouts/MainLayout.vue'
import { pojectsRoute } from './projects'
import { membersRoute } from './members'

export const MAIN_PATH = '/main'
export const MAIN_PATH_NAME = 'main'

export const mainRoute: RouteRecordRaw = {
  name: MAIN_PATH_NAME,
  path: MAIN_PATH,
  redirect: '/projects',
  component: MainLayout,
  children: [pojectsRoute, membersRoute],
}
