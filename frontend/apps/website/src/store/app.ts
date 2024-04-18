import { defineStore } from 'pinia'
import { pinia } from '@/plugins'
import { apiGetGithubOAuthConfig, oAuthURLMap } from '@/api/sign/oAuth'

interface AppState {
  // 全局loading 计数器
  globalLoadingIndicator: number
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
    globalLoadingIndicator: 0,
    oAuthPlatformConfig: {
      github: null,
    },
  }),
  getters: {
    isShowGlobalLoading: state => state.globalLoadingIndicator > 0,
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
    },
    hideGlobalLoading() {
      this.globalLoadingIndicator && this.globalLoadingIndicator--
    },
  },
})

export default useAppStore

export const useAppStoreWithOut = () => useAppStore(pinia)
