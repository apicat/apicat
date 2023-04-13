import type { RouteRecordRaw } from 'vue-router'
import ProjectListPage from '@/views/project/ProjectListPage.vue'

/**
 * project module routes
 */
export const pojectsRoute: RouteRecordRaw = {
  name: 'projects',
  path: '/projects',
  alias: '/home',
  component: ProjectListPage,
  meta: {
    title: 'app.project.list.title',
  },
}
