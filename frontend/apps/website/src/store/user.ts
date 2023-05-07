import { defineStore } from 'pinia'
import { router } from '@/router'
import { MAIN_PATH, LOGIN_PATH } from '@/router'
import Storage from '@/commons/storage'
import { userEmailLogin, userRegister } from '@/api/user'
import { UserInfo } from '@/typings/user'

interface UserState {
  userInfo: {}
  token: string | null
}

export const useUserStore = defineStore({
  id: 'user',

  state: (): UserState => ({
    token: Storage.get(Storage.KEYS.TOKEN) || null,
    userInfo: Storage.get(Storage.KEYS.USER) || null,
  }),

  getters: {
    isLogin: (state) => !!state.token,
  },

  actions: {
    updateToken(token: string) {
      Storage.set(Storage.KEYS.TOKEN, token)
      this.token = token
    },

    // 登录
    async login(form: any) {
      try {
        const res = await userEmailLogin(form)
        await this.getUserInfo()
        this.goHome()
        return res
      } catch (error) {
        //
      }
    },

    async register(form: UserInfo) {
      try {
        const res = await userRegister(form)
        this.goHome()
        return res
      } catch (error) {
        //
      }
    },
    // 退出
    async logout() {
      Storage.removeAll([Storage.KEYS.TOKEN, Storage.KEYS.USER])
      try {
        // await logout()
      } catch (error) {
        //
      } finally {
        this.token = null
        this.userInfo = {} as any
        location.href = LOGIN_PATH
      }
    },

    // 获取个人信息
    async getUserInfo() {
      // try {
      //   const { data } = (await getUserInfo()) || {}
      //   data.address = [data.province, data.city]
      //   this.userInfo = data
      //   Storage.set(Storage.KEYS.USER_INFO, this.userInfo)
      //   return this.userInfo
      // } catch (error) {
      //   //
      // }
      // return {}
    },

    goHome(path?: string) {
      router.replace(path || MAIN_PATH)
    },

    // 更新个人信息
    async updateUserInfo(user: any) {
      // await updateUserProfile(user)
      // this.$patch({ userInfo: { ...user } })
      // Storage.set(Storage.KEYS.USER_INFO, this.userInfo)
    },
  },
})
