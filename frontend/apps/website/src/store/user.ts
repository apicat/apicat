import { defineStore } from 'pinia'
import useLocaleStore from './locale'
import { LOGIN_PATH, MAIN_PATH, router } from '@/router'
import Storage from '@/commons/storage'
import { apiLogin, apiRegister } from '@/api/sign/user'
import { pinia } from '@/plugins'
import { apiGetUserInfo } from '@/api/user'

interface UserState {
  userInfo: UserAPI.ResponseUserInfo | Record<string, any>
  token: string | null
}

export const useUserStore = defineStore({
  id: 'user',

  state: (): UserState => ({
    token: Storage.get(Storage.KEYS.TOKEN) || null,
    userInfo: {},
  }),

  getters: {
    isLogin: state => !!state.token,
    isAdmin: state => state.userInfo.role === 'admin',
  },

  actions: {
    // 登录
    async login(form: SignAPI.RequestLogin, url?: string) {
      try {
        const data = await apiLogin(form)
        await this.afterSign(data, url)
        return data
      }
      catch (error) {
        //
      }
    },

    async afterSign(data: SignAPI.ResponseLogin | SignAPI.ResponseRegister, url?: string) {
      this.updateToken(data.accessToken)
      this.getUserInfo()
      this.goHome(url)
    },

    // 注册
    async register(form: SignAPI.RequestRegister, url?: string) {
      const data = await apiRegister(form)
      await this.afterSign(data, url)
      return data
    },

    async getUserInfo(): Promise<UserAPI.ResponseUserInfo | any> {
      try {
        const user = await apiGetUserInfo()
        this.updateUserInfo(user)
      }
      catch (error) {
        //
      }
      return this.userInfo
    },

    // 退出
    logout() {
      Storage.removeAll([Storage.KEYS.TOKEN, Storage.KEYS.USER, Storage.KEYS.SELECTED_PROJECT_GROUP])
      this.token = null
      this.userInfo = {} as any
      location.href = LOGIN_PATH
    },

    updateToken(token: string) {
      if (!token)
        return

      Storage.set(Storage.KEYS.TOKEN, token)
      this.token = token
    },

    goHome(path?: string) {
      router.replace(path || MAIN_PATH)
    },

    // 更新个人信息
    async updateUserInfo(user: Partial<UserAPI.ResponseUserInfo>) {
      Object.assign(this.userInfo, user)
      const { switchLanguage } = useLocaleStore()
      user.language && await switchLanguage(user.language)
    },
  },
})

export const useUserStoreWithOut = () => useUserStore(pinia)
