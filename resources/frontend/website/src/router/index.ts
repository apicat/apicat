import { createRouter, createWebHistory } from 'vue-router'
import { MAIN_PATH } from './constant'
import MainLayout from '../layout/MainLayout.vue'

// 各模块路由
import { indexRouter, loginRouter, registerRouter, notFoundRouter } from './base.router'
import userRouters from './user.router'
import membersRouters from './members.router'
import projectsRouters from './projects.router'
import projectRouters from './project.router'
import documentRouter from './document.router'
import iterationsRouter from './iterations.router'
import iterateDocumentRouter from './iterate.router'

import { documentPreviewRouters, projectPreviewRouters, trashPreviewRouters } from './preview.router'

export const router = createRouter({
    history: createWebHistory(import.meta.env.BASE_URL),
    routes: [
        indexRouter,
        loginRouter,
        registerRouter,

        documentRouter,

        iterateDocumentRouter,

        {
            path: MAIN_PATH,
            name: 'main',
            component: MainLayout,
            redirect: { name: 'projects' },
            children: [projectsRouters, userRouters, membersRouters, projectRouters, iterationsRouter],
        },

        projectPreviewRouters,
        documentPreviewRouters,
        trashPreviewRouters,

        notFoundRouter,
    ],
})

export default router

export * from './filter'
