import { defineStore } from 'pinia'
import { router } from '@/router'
import { MAIN_PATH, LOGIN_PATH } from '@/router'
import Storage from '@/commons/storage'
import { userEmailLogin, userRegister, modifyUserInfo, modifyPassword, getUserInfo } from '@/api/user'
import { UserInfo, UserRoleInTeam, UserRoleInTeamMap } from '@/typings/user'
import { pinia } from '@/plugins'

interface UserState {
  userInfo: UserInfo | Record<string, any>
  token: string | null
}

export const useUserStore = defineStore({
  id: 'user',

  state: (): UserState => ({
    token: Storage.get(Storage.KEYS.TOKEN) || null,
    userInfo: {},
  }),

  getters: {
    isLogin: (state) => !!state.token,
    isSuperAdmin: (state) => state.userInfo?.role === UserRoleInTeam.SUPER_ADMIN,
    isNormalUser: (state) => state.userInfo?.role === UserRoleInTeam.USER,
    userRoles: () =>
      Object.keys(UserRoleInTeamMap)
        .filter((key: string) => key !== UserRoleInTeam.SUPER_ADMIN)
        .map((key: string) => {
          return {
            text: (UserRoleInTeamMap as any)[key],
            value: key,
          }
        }),
  },

  actions: {
    // 登录
    async login(form: any) {
      try {
        const data: any = await userEmailLogin(form)
        this.updateToken(data.access_token)
        this.updateUserInfo(data.user)
        this.goHome()
        return data
      } catch (error) {
        //
      }
    },

    async register(form: UserInfo) {
      try {
        const data: any = await userRegister(form)
        this.updateToken(data.access_token)
        this.updateUserInfo(data.user)
        this.goHome()
        return data
      } catch (error) {
        //
      }
    },

    async getUserInfo(): Promise<UserInfo | void> {
      try {
        const user: any = await getUserInfo()
        this.updateUserInfo(user)
        return user
      } catch (error) {
        //
      }
    },

    async modifyUserInfo(form: UserInfo) {
      try {
        const user: any = await modifyUserInfo(form)
        this.updateUserInfo(user)
      } catch (error) {
        //
      }
    },

    async modifyUserPassword(form: UserInfo) {
      try {
        await modifyPassword(form)
      } catch (error) {
        //
      }
    },
    // 退出
    logout() {
      Storage.removeAll([Storage.KEYS.TOKEN, Storage.KEYS.USER])
      this.token = null
      this.userInfo = {} as any
      this.goHome(LOGIN_PATH)
    },

    updateToken(token: string) {
      Storage.set(Storage.KEYS.TOKEN, token)
      this.token = token
    },

    goHome(path?: string) {
      router.replace(path || MAIN_PATH)
    },

    // 更新个人信息
    updateUserInfo(user: UserInfo) {
      this.$patch({ userInfo: { ...this.userInfo, ...user } })
    },
  },
})

export const useUserStoreWithOut = () => useUserStore(pinia)
