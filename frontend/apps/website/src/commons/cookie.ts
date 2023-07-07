import { getCookie, setCookie, removeCookie } from 'tiny-cookie'
import { STORAGE_PREFIX } from './constant'

export interface CookieOptions {
  domain?: string
  path?: string
  expires?: Date | string | number
  'max-age'?: number
  secure?: boolean
  samesite?: string
}

export const Cookies = {
  KEYS: {
    SHARE_PROJECT: `${STORAGE_PREFIX}.share.p.`,
    SHARE_DOCUMENT: `${STORAGE_PREFIX}.share.d.`,
  },
  get: (key: string) => {
    try {
      return getCookie(key, JSON.parse)
    } catch (error) {
      return null
    }
  },
  set: <T>(key: string, value: T, options?: CookieOptions) => setCookie(key, value, JSON.stringify, options),
  remove: removeCookie,
}
