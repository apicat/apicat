import { RouteRecordRaw } from 'vue-router'
import MainLayout from '@/layouts/MainLayout.vue'
import { pojectsRoute } from './projects'
import { membersRoute } from './members'
import { MAIN_PATH, MAIN_PATH_NAME } from './constant'

export const mainRoute: RouteRecordRaw = {
  name: MAIN_PATH_NAME,
  path: MAIN_PATH,
  redirect: '/projects',
  component: MainLayout,
  children: [pojectsRoute, membersRoute],
}
