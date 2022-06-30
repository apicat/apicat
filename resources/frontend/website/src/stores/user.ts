import { userEmailLogin, userRegister, getUserInfo, logout, updateUserProfile } from '@/api/user'
import { getFirstChar, Storage } from '@natosoft/shared'
import { TEAM_ROLE } from '@/common/constant'
import { defineStore } from 'pinia'
import { UserInfo } from '~/store'
import { router } from '@/router'
import { MAIN_PATH, LOGIN_PATH } from '@/router/constant'
import { isEmpty } from 'lodash'

interface UserState {
    userInfo: UserInfo
    token: string | null
}

export const useUserStore = defineStore({
    id: 'user',

    state: (): UserState => ({
        token: Storage.get(Storage.KEYS.TOKEN) || null,
        userInfo: Storage.get(Storage.KEYS.USER_INFO) || {},
    }),

    getters: {
        isLogin: (state) => !!state.token,
        isAdmin: (state) => state.userInfo.authority === TEAM_ROLE.ADMIN,
        isManager: (state) => state.userInfo.authority === TEAM_ROLE.MANAGER,
        isNormal: (state) => state.userInfo.authority === TEAM_ROLE.NORMAL,
        lastName: (state) => getFirstChar(state.userInfo.name),
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

        async register(form: any) {
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
            Storage.removeAll([Storage.KEYS.TOKEN, Storage.KEYS.USER_INFO, Storage.KEYS.ACTIVE_PROJECT_GROUP])
            try {
                await logout()
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
            if (!isEmpty(this.userInfo)) {
                return this.userInfo
            }

            try {
                const { data } = (await getUserInfo()) || {}
                Storage.set(Storage.KEYS.USER_INFO, data)
                this.userInfo = data

                return this.userInfo
            } catch (error) {
                //
            }

            return {}
        },

        goHome(path?: string) {
            router.replace(path || MAIN_PATH)
        },

        // 更新个人信息
        async updateUserInfo(user: any) {
            await updateUserProfile(user)
            this.$patch({ userInfo: { name: user.name, email: user.email, avatar: user.avatar } })
            Storage.set(Storage.KEYS.USER_INFO, this.userInfo)
        },
    },
})
