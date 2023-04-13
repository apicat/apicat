import { STORAGE_PREFIX } from './constant'

interface StorageKeys {
  TOKEN: string
  USER_INFO: string
  LOCALE: string
}

const Storage = {
  KEYS: {
    TOKEN: `${STORAGE_PREFIX}.token`,
    LOCALE: `${STORAGE_PREFIX}.locale`,
  } as StorageKeys,

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

  removeAll(keys = []) {
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

Storage.get(Storage.KEYS.TOKEN)

export default Storage
