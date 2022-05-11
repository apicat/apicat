import { RouteRecordRaw } from 'vue-router'

import ProjectsNav from '@/views/projects/ProjectsNav.vue'

const Projects = () => import('../views/projects/Projects.vue')

export const ProjectListRoute = { path: '', name: 'projects.index', meta: { title: '项目列表' }, component: Projects }

export const ProjectRoutes: RouteRecordRaw[] = [ProjectListRoute]

export default {
    path: '/projects',
    name: 'projects',
    redirect: { name: 'projects.index' },
    component: ProjectsNav,
    children: ProjectRoutes,
}
