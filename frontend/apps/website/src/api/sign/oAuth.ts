import DefaultAjax from '../Ajax'
import { OAuthPlatform, OAuthPlatformConfig } from '@/commons/constant'
import { OAUTH_CONNECT_NAME, OAUTH_NAME } from '@/router/constant'
import router from '@/router'

export function createOAuthLoginCallbackURL(type: OAuthPlatform) {
  const base = window.location.origin
  const path = router.resolve({
    name: OAUTH_NAME,
    params: {
      type,
    },
  })
  return base + path.fullPath
}

export function createOAuthConnectCallbackURL(type: OAuthPlatform) {
  const base = window.location.origin
  const path = router.resolve({
    name: OAUTH_CONNECT_NAME,
    params: {
      type,
    },
  })
  return base + path.fullPath
}

export const oAuthURLMap = {
  github: (redirectUrl?: string, params?: Record<string, any>) => {
    params = new URLSearchParams({
      redirect_uri: redirectUrl || createOAuthLoginCallbackURL(OAuthPlatform.GITHUB),
      ...(params || {}),
      ...OAuthPlatformConfig.GITHUB.params,
    })
    return `${OAuthPlatformConfig.GITHUB.OAUTH_URL}?${params.toString()}`
  },
}

export async function apiOAuthLoginWithCode(platform: string, data: SignAPI.RequestOAuthLogin): Promise<SignAPI.ResponseOAuthLogin | null> {
  switch (platform) {
    case OAuthPlatform.GITHUB:
      return DefaultAjax.post(`/account/oauth/${platform}/login`, data)

    default:
      return null
  }
}

// 获取GitHub OAuth配置
export async function apiGetGithubOAuthConfig(): Promise<{ clientID: string }> {
  return DefaultAjax.get('/sysconfigs/github',{},{ isShowErrorMsg: false })
}
