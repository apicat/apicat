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
import { NoTeamError } from '@/api/error'

const whiteList = [TEAM_CREATE_NAME, TEAM_JOIN_NAME, USER_PAGE_NAME, OAUTH_NAME, USER_SETTING_NAME, OAUTH_CONNECT_NAME, SYSTEM_SETTING_NAME, SYSTEM_PAGE_NAME]

export function setupGetTeamFilter(router: Router) {
  router.beforeEach(async (to, __, next) => {
    // 默认拦截所有需要登录的页面
    if (!to.meta.ignoreAuth) {
      const teamStore = useTeamStore()
      let err
      try {
        // 未init获取
        if (!teamStore.inited)
          await teamStore.init()
      }
      catch (e) {
        err = e
      }

      // 如果在whitelist中直接走
      if (whiteList.includes(to.name as string))
        return next()
      else if (err instanceof NoTeamError)
        return next(TEAM_CREATE_PATH)

      return next()
    }

    next()
  })
}
