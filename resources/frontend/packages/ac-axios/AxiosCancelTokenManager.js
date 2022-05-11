import axios from 'axios'
import { isFunction } from '@ac/shared'

export default class AxiosPendingManager {
    constructor() {
        this._pendingMap = new Map()
    }

    gatherKey(config) {
        return [config.url, config.method, JSON.stringify(config.params)].join('_')
    }

    add(config) {
        this.remove(config)

        const tokenSource = axios.CancelToken.source()

        config.cancelToken = config.cancelToken || tokenSource.token

        const key = this.gatherKey(config)
        this._pendingMap.set(key, tokenSource.cancel)
    }

    remove(config) {
        const key = this.gatherKey(config)
        if (this._pendingMap.has(key)) {
            const cancel = this._pendingMap.get(key)
            cancel && isFunction(cancel) && cancel()
            this._pendingMap.delete(key)
        }
    }

    removeAll() {
        this._pendingMap.forEach((cancel) => cancel && isFunction(cancel) && cancel())
        this.clear()
    }

    clear() {
        this._pendingMap.clear()
    }
}
