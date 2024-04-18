import type { Router } from 'vue-router'
import {
  OAUTH_CONNECT_NAME,
  OAUTH_NAME,
  SYSTEM_PAGE_NAME,
  SYSTEM_SETTING_NAME,
  TEAM_CREATE_NAME,
  TEAM_CREATE_PATH,
  TEAM_JOIN_NAME,
  USER_PAGE_NAME,
  USER_SETTING_NAME,
} from '../constant'
import { useTeamStore } from '@/store/team'

const whiteList = [TEAM_CREATE_NAME, TEAM_JOIN_NAME, USER_PAGE_NAME, OAUTH_NAME, USER_SETTING_NAME, OAUTH_CONNECT_NAME, SYSTEM_SETTING_NAME, SYSTEM_PAGE_NAME]

export function setupGetTeamFilter(router: Router) {
  let tried = false
  router.beforeEach(async (to, __, next) => {
    // 默认拦截所有需要登录的页面
    if (!to.meta.ignoreAuth) {
      const teamStore = useTeamStore()
      if (!tried) {
        await teamStore.init()
        tried = true
      }

      // 白名单
      if (whiteList.includes(to.name as string))
        return next()

      // 获取团队信息后，仍然没有激活团队，跳转到创建团队页面
      if (!teamStore.currentTeam)
        return next(TEAM_CREATE_PATH)

      return next()
    }

    next()
  })
}
