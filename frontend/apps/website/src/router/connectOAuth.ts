import type { NavigationGuardNext, RouteLocationNormalized, RouteRecordRaw } from 'vue-router'
import {
  COMPLETE_INFO_PATH,
  NOT_FOUND_PATH,
  OAUTH_CONNECT_NAME,
  OAUTH_CONNECT_PATH,
  USER_PAGE_NAME,
} from '@/router/constant'
import { apiConnectOAuth } from '@/api/user'
import type { OAuthPlatform } from '@/commons/constant'

export const OAUTH_PLATFORMS: Record<OAuthPlatform, string> = {
  github: COMPLETE_INFO_PATH,
}

export const connectOAuthRoute: RouteRecordRaw = {
  path: OAUTH_CONNECT_PATH,
  name: OAUTH_CONNECT_NAME,
  meta: { title: 'app.pageTitles.connectOAuth' },
  component: { template: '' },
  beforeEnter: async (to: RouteLocationNormalized, _: RouteLocationNormalized, next: NavigationGuardNext) => {
    const platform = to.params.type as OAuthPlatform
    if (!platform || !OAUTH_PLATFORMS[platform])
      return next(NOT_FOUND_PATH)

    try {
      await apiConnectOAuth(platform as unknown as SignAPI.OAuthPlatform, { code: to.query.code as string })
    }
    catch (error) {
      //
    }

    return next({
      name: USER_PAGE_NAME,
      params: {
        page: platform,
      },
    })
  },
}
