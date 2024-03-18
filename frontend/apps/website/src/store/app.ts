import { defineStore } from 'pinia'
import { pinia } from '@/plugins'
import { apiGetGithubOAuthConfig, oAuthURLMap } from '@/api/sign/oAuth'

interface AppState {
  globalLoadingIndicator: number
  isShowGlobalLoading: boolean
  oAuthPlatformConfig: {
    github: null | {
      client_id: string
    }
  }
}

interface OAuthURLConfig {
  github(redirectUrl?: string): string
}

export const useAppStore = defineStore('app', {
  state: (): AppState => ({
    globalLoadingIndicator: 1,
    isShowGlobalLoading: false,
    oAuthPlatformConfig: {
      github: null,
    },
  }),
  getters: {
    isShowGithubOAuth: state => state.oAuthPlatformConfig.github !== null,
    oAuthURLConfig: (state): OAuthURLConfig => {
      return {
        github: (redirectUrl?: string) => oAuthURLMap.github(redirectUrl, state.oAuthPlatformConfig.github || {}),
      }
    },
  },
  actions: {
    async initAppConfig() {
      await this.getGithubOAuthConfig()
    },

    // 获取github oauth 配置
    async getGithubOAuthConfig() {
      try {
        const { clientID } = await apiGetGithubOAuthConfig()
        this.updatGithuClienId(clientID)
      }
      catch (error) {
        //
      }
    },

    // 更新github oauth client_id
    updatGithuClienId(client_id: string) {
      if (!client_id)
        return

      this.oAuthPlatformConfig.github = {
        client_id,
      }
    },

    showGlobalLoading() {
      this.globalLoadingIndicator++
      this.isShowGlobalLoading = true
    },

    hideGlobalLoading() {
      this.globalLoadingIndicator--
      this.isShowGlobalLoading = false
    },
  },
})

export default useAppStore

export const useAppStoreWithOut = () => useAppStore(pinia)
