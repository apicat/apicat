import { RouteRecordRaw } from 'vue-router'
import ProjectsNav from '@/views/projects/ProjectsNav.vue'
import { MAIN_PATH } from '@/router/constant'

const Projects = () => import('../views/projects/Projects.vue')

export const ProjectListRoute = { path: '', name: 'projects.index', meta: { title: '项目列表' }, component: Projects }

export const ProjectRoutes: RouteRecordRaw[] = [ProjectListRoute]

export const ProjectRootRoute = {
    path: MAIN_PATH,
    name: 'projects',
    meta: { title: '项目', activeClassPrefixes: 'project', browserTitle: '项目 - ApiCat' },
    redirect: { name: 'projects.index' },
}

export default {
    ...ProjectRootRoute,
    component: ProjectsNav,
    children: ProjectRoutes,
}
