import { createRouter, createWebHistory } from 'vue-router'
import { MAIN_PATH } from './constant'
import MainLayout from '../layout/MainLayout.vue'

// 各模块路由
import { indexRouter, loginRouter, registerRouter, notFoundRouter } from './base.router'
import UserRouters from './user.router'
import MembersRouters from './members.router'
import ProjectsRouters from './projects.router'
import ProjectRouters from './project.router'
import DocumentRouter from './document.router'
import { documentPreviewRouters, projectPreviewRouters, trashPreviewRouters } from './preview.router'

export const router = createRouter({
    history: createWebHistory(import.meta.env.BASE_URL),
    routes: [
        indexRouter,
        loginRouter,
        registerRouter,

        {
            path: MAIN_PATH,
            name: 'main',
            component: MainLayout,
            redirect: { name: 'projects' },
            children: [ProjectsRouters, UserRouters, MembersRouters, ProjectRouters],
        },

        DocumentRouter,
        projectPreviewRouters,
        documentPreviewRouters,
        trashPreviewRouters,

        notFoundRouter,
    ],
})

export default router

export * from './filter'
