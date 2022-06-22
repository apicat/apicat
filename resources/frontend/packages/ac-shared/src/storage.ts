/* eslint-disable no-global-assign */
/**
 * @file Storage
 * @brief 本地存储工具
 */

const KEY_PREFIX = 'api.cat'

export interface IStorage {
    KEYS: {
        TOKEN: string
        USER_INFO: string
        ACTIVE_PROJECT_GROUP: string
        SECRET_DOCUMENT_TOKEN: string
        SECRET_PROJECT_TOKEN: string
    }
    getStorage(isSession: boolean | undefined): Storage
    get(key: string, isSession?: boolean): any | undefined
    set(key: string, value: unknown, isSession?: boolean): void
    isLocalStorage(): boolean
    remove(key: string, isSession?: boolean): undefined | IStorage
    removeAll(keys: Array<string>): void
}

export const Storage: IStorage = {
    KEYS: {
        TOKEN: `${KEY_PREFIX}.token`,
        USER_INFO: `${KEY_PREFIX}.user.info`,
        ACTIVE_PROJECT_GROUP: `${KEY_PREFIX}.active.project.group`,
        SECRET_DOCUMENT_TOKEN: `${KEY_PREFIX}.preview.document.`,
        SECRET_PROJECT_TOKEN: `${KEY_PREFIX}.preview.project.`,
    },

    get(key: string, isSession?: boolean) {
        if (!this.isLocalStorage()) {
            return undefined
        }
        const value = this.getStorage(isSession).getItem(key)
        if (value) return JSON.parse(value)
        return undefined
    },

    set(key: string, value: any, isSession?: boolean) {
        if (!this.isLocalStorage()) {
            return
        }
        value = JSON.stringify(value)
        this.getStorage(isSession).setItem(key, value)

        return
    },

    remove(key: string, isSession?: boolean): undefined | IStorage {
        if (!this.isLocalStorage()) return undefined
        this.getStorage(isSession).removeItem(key)
        return this
    },

    removeAll(keys = []) {
        ;(keys.length ? keys : Object.keys(this.KEYS)).forEach((item) => this.remove(keys.length ? item : (this.KEYS as any)[item]))
    },

    getStorage(isSession?: undefined | boolean): Storage {
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
