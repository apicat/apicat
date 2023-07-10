import { STORAGE_PREFIX } from './constant'

const Storage = {
  KEYS: {
    TOKEN: `${STORAGE_PREFIX}.token`,
    USER: `${STORAGE_PREFIX}.user`,
    LOCALE: `${STORAGE_PREFIX}.locale`,
    CODE_GENERATE_CONFIG: `${STORAGE_PREFIX}.code.generate`,
    CODE_GENERATE_LANGUAGE: `${STORAGE_PREFIX}.code.gen.lang`,
  },

  get(key: string, isSession?: boolean) {
    if (!this.isLocalStorage()) {
      return null
    }
    const value = this.getStorage(isSession).getItem(key)
    if (value) {
      try {
        return JSON.parse(value)
      } catch (error) {
        return null
      }
    }
    return null
  },

  set(key: string, value: unknown, isSession?: boolean) {
    if (!this.isLocalStorage()) {
      return null
    }
    value = JSON.stringify(value)
    this.getStorage(isSession).setItem(key, value as string)

    return this
  },

  remove(key: string, isSession?: boolean) {
    if (!this.isLocalStorage()) return null
    this.getStorage(isSession).removeItem(key)
    return this
  },

  removeAll(keys: string[]) {
    ;(keys.length ? keys : Object.keys(this.KEYS)).forEach((key: string) => this.remove(key))
  },

  getStorage(isSession: boolean = false) {
    return isSession ? sessionStorage : localStorage
  },

  isLocalStorage() {
    try {
      if (!window.localStorage) {
        return false
      }
      return true
    } catch (e) {
      return false
    }
  },
}

export default Storage
