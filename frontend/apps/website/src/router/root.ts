import type { RouteRecordRaw } from 'vue-router'
import { MAIN_PATH, ROOT_PATH, ROOT_PATH_NAME } from './constant'

export const rootRoute: RouteRecordRaw = {
  name: ROOT_PATH_NAME,
  path: ROOT_PATH,
  redirect: MAIN_PATH,
  meta: { ignoreAuth: true },
}
