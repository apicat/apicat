import type { RouteRecordRaw } from 'vue-router'

const ProjectListPage = () => import('@/views/project/ProjectListPage.vue')

export const pojectsRoute: RouteRecordRaw = {
  name: 'projects',
  path: '/projects',
  alias: '/home',
  component: ProjectListPage,
  meta: {},
}
