import type { NavigationGuardNext, RouteLocationNormalized, RouteRecordRaw } from 'vue-router'
import { apiOAuthLoginWithCode } from '@/api/sign/oAuth'
import {
  COMPLETE_INFO_NAME,
  COMPLETE_INFO_PATH,
  LOGIN_PATH,
  MAIN_PATH,
  NOT_FOUND_PATH,
  OAUTH_NAME,
  OAUTH_PATH,
} from '@/router/constant'
import { useUserStore } from '@/store/user'
import type { OAuthPlatform } from '@/commons/constant'
import { DEFAULT_LANGUAGE } from '@/commons/constant'
import { flattenObject } from '@/commons'
import { useGlobalLoading } from '@/hooks/useGlobalLoading'

export const OAUTH_PLATFORMS: Record<OAuthPlatform, string> = {
  github: COMPLETE_INFO_PATH,
}

export const oauthRoute: RouteRecordRaw = {
  path: OAUTH_PATH,
  name: OAUTH_NAME,
  meta: { ignoreAuth: true, title: 'OAuth Login' },
  component: { template: '' },
  beforeEnter: async (to: RouteLocationNormalized, _: RouteLocationNormalized, next: NavigationGuardNext) => {
    const { showGlobalLoading, hideGlobalLoading } = useGlobalLoading()
    showGlobalLoading()
    try {
      const platform = to.params.type as OAuthPlatform
      if (!platform || !OAUTH_PLATFORMS[platform])
        return next(NOT_FOUND_PATH)

      const res = await apiOAuthLoginWithCode(platform, {
        code: to.query.code as string,
        invitationToken: to.query.invitationToken as string,
        language: DEFAULT_LANGUAGE,
      })
      if (!res)
        return next(NOT_FOUND_PATH)

      if (res.accessToken) {
        useUserStore().updateToken(res.accessToken)
        return next(MAIN_PATH)
      }
      return next({
        name: COMPLETE_INFO_NAME,
        params: {
          ...to.params,
        },
        query: {
          ...to.query,
          ...flattenObject(res),
        },
      })
    }
    catch (err) {
      return next(LOGIN_PATH)
    }
    finally {
      hideGlobalLoading()
    }
  },
}
