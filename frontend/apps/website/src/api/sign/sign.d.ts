declare namespace SignAPI {
  type Languages = typeof import('@/commons/constant').Languages
  type OAuthPlatform = typeof import('@/commons/constant').OAuthPlatform
  interface ResponseLogin {
    accessToken: string
  }
  interface RequestLogin {
    email: string
    password: string
    invitationToken?: string
  }

  interface RegisterUserBind {
    oauthUserID: string
    type: string
  }
  interface RequestRegister {
    email: string
    name: string
    password: string
    avatar?: string
    bind?: RegisterUserBind
    invitationToken?: string
    level?: number
    language?: keyof Languages
  }
  interface ResponseRegister {
    accessToken: string
  }

  interface ResponseOAuthLogin {
    email: string
    name: string
    oauthUserID: string
    type: string
    accessToken?: string
    avatar?: string
    level?: number
  }

  interface RequestOAuthLogin {
    code: string
    language: keyof Languages
    invitationToken?: string
  }
}
