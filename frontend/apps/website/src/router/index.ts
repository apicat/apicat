import { createRouter, createWebHistory } from 'vue-router'
import { rootRoute } from './root'
import { mainRoute } from './main'
import { noPermissionRoute, notFoundRoute } from './base'
import {
  completeInfoRoute,
  forgetPassRoute,
  loginRoute,
  registerRoute,
  registerVerificationEmailRoute,
  resetPassRoute,
  verificationEmailForModifyRoute,
} from './sign'
import { oauthRoute } from './oauth'
import { connectOAuthRoute } from './connectOAuth'
import { collectionShareRoute } from './collection'
import { collectionHistoryRoute, schemaHistoryRoute } from './history'
import { iterationDetailRoute } from './iterationDetail'
import { joinTeamRoute } from '@/router/team'
import { projectDetailRoute } from '@/router/projectDetail'

export const router = createRouter({
  history: createWebHistory(),
  routes: [
    rootRoute,
    mainRoute,

    // sign
    loginRoute,
    registerRoute,
    forgetPassRoute,
    resetPassRoute,
    registerVerificationEmailRoute,
    verificationEmailForModifyRoute,

    // oauth
    oauthRoute,
    connectOAuthRoute,
    completeInfoRoute,

    // team
    joinTeamRoute,

    // project
    projectDetailRoute,

    // schema
    schemaHistoryRoute,

    // collcetion
    collectionShareRoute,
    collectionHistoryRoute,

    // iteration
    iterationDetailRoute,

    noPermissionRoute,
    notFoundRoute,
  ],
})

export * from './oauth'
export * from './base'
export * from './sign'
export * from './root'
export * from './main'

// 路由拦截器
export * from './filter'
export * from './constant'

export default router
