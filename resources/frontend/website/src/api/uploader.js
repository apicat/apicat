import { Storage } from '@ac/shared'
import { uploadPath } from '@/api/upload'
import { noop } from 'lodash-es'

function Uploader() {
    // 1MB
    this.chunkSize = 1024 * 1024
    this.threadsQuantity = 2
    this.timeout = 1000 * 10

    // size -> bit

    this.file = null
    this.aborted = false
    this.uploadedSize = 0
    this.progressCache = {}
    this.activeConnections = {}
    this.taskConfig = {}
}

Uploader.prototype.setOptions = function (options = {}) {
    this.chunkSize = options.chunkSize || this.chunkSize
    this.threadsQuantity = options.threadsQuantity || this.threadsQuantity
}

Uploader.prototype.setupFile = function (file) {
    if (!file) {
        return
    }

    this.file = file
}

Uploader.prototype.start = function () {
    if (!this.file) {
        throw new Error("Can't start uploading: file have not chosen")
    }

    const chunksQuantity = Math.ceil(this.file.size / this.chunkSize)
    this.chunksQueue = new Array(chunksQuantity)
        .fill()
        .map((_, index) => index)
        .reverse()

    const xhr = new XMLHttpRequest()
    xhr.timeout = this.timeout
    xhr.open('post', '/api/upload_init')
    xhr.setRequestHeader('Content-Type', 'application/json;charset=UTF-8')
    xhr.setRequestHeader('Accept', 'application/json, text/plain, */*')
    let authorization = Storage.get(Storage.KEYS.TOKEN) || ''

    if (authorization) {
        xhr.setRequestHeader('Authorization', authorization)
    }

    xhr.withCredentials = true

    xhr.onreadystatechange = () => {
        if (xhr.readyState === XMLHttpRequest.DONE && xhr.status === 200) {
            const { status, msg, data } = JSON.parse(xhr.responseText)
            if (status !== 0) {
                this.complete(new Error(msg || '上传初始化失败！'))
                return
            }

            this.fileId = data.file_id
            this.sendNext()
        }
    }

    xhr.onerror = (error) => {
        this.complete(error)
    }

    xhr.ontimeout = (error) => {
        this.complete(error)
    }

    xhr.send(
        JSON.stringify({
            name: this.file.name,
            chunks: chunksQuantity,
            fileSize: Math.ceil(this.file.size / 1024),
        })
    )
}

Uploader.prototype.sendNext = function () {
    const activeConnections = Object.keys(this.activeConnections).length

    // 某一任务出现三次失败后 后续不在上传；
    let isBreak = false
    Object.keys(this.taskConfig).forEach((key) => {
        if (this.taskConfig[key].retryCount >= 3) {
            isBreak = true
        }
    })

    if (isBreak) {
        this.complete(new Error('上传失败！'))
        return
    }

    if (activeConnections >= this.threadsQuantity) {
        return
    }

    // 上传完毕
    if (!this.chunksQueue.length) {
        if (!activeConnections) {
            this.getResult()
        }

        return
    }

    const chunkId = this.chunksQueue.pop()
    const sentSize = chunkId * this.chunkSize
    const chunk = this.file.slice(sentSize, sentSize + this.chunkSize)

    this.sendChunk(chunk, chunkId)
        .then(() => this.sendNext())
        .catch((error) => {
            this.chunksQueue.push(chunkId)
            let config = this.taskConfig[chunkId]
            if (config && config.retryCount <= 3) {
                // console.log(`尝试第 ${config.retryCount} 次重新上传：${chunkId}`);
                this.sendNext()
            } else {
                this.complete(error)
            }
        })

    this.sendNext()
}

Uploader.prototype.complete = function (error) {
    const message = (error && error.message) || error

    if (error && !this.aborted) {
        this.end(message)
        return
    }

    this.end(message)
}

Uploader.prototype.sendChunk = function (chunk, id) {
    return this.upload(id, chunk)
}

Uploader.prototype.handleProgress = function (chunkId, event) {
    if (event.type === 'progress' || event.type === 'error' || event.type === 'abort') {
        this.progressCache[chunkId] = event.loaded
    }

    if (event.type === 'loadend') {
        this.uploadedSize += this.progressCache[chunkId] || 0
        delete this.progressCache[chunkId]
    }

    const inProgress = Object.keys(this.progressCache).reduce((memo, id) => (memo += this.progressCache[id]), 0)

    const sendedLength = Math.min(this.uploadedSize + inProgress, this.file.size)

    this.onProgress &&
        this.onProgress({
            loaded: sendedLength,
            total: this.file.size,
        })
}

Uploader.prototype.upload = function (id, file) {
    return new Promise((resolve, reject) => {
        let taskConfig = (this.taskConfig[id] = this.taskConfig[id] || { retryCount: 0 })

        const xhr = (this.activeConnections[id] = new XMLHttpRequest())
        xhr.timeout = this.timeout

        const progressListener = this.handleProgress.bind(this, id)

        xhr.upload.addEventListener('progress', progressListener)

        xhr.addEventListener('error', progressListener)
        xhr.addEventListener('abort', progressListener)
        xhr.addEventListener('loadend', progressListener)

        xhr.open('post', '/api/upload_save')

        xhr.setRequestHeader('Accept', 'application/json, text/plain, */*')
        let authorization = Storage.get(Storage.KEYS.TOKEN) || ''
        if (authorization) {
            xhr.setRequestHeader('Authorization', authorization)
        }

        xhr.onreadystatechange = () => {
            delete this.activeConnections[id]

            if (xhr.readyState === XMLHttpRequest.DONE && xhr.status === 200) {
                const req = JSON.parse(xhr.responseText)

                if (req.status !== 0) {
                    reject(new Error(req.msg || '上传失败！'))
                    taskConfig.retryCount++
                    return
                }

                if (req.status === 0) {
                    resolve(req)
                }
            }

            // 500 other 不需要重新尝试
            if (xhr.readyState === XMLHttpRequest.DONE && xhr.status !== 200) {
                reject(new Error('上传失败！'))
                taskConfig.retryCount++
            }
        }

        // 不需要重新尝试
        xhr.onerror = (error) => {
            reject(error)
            delete this.activeConnections[id]
        }

        // 取消上传
        xhr.onabort = () => {
            reject(new Error('取消上传！'))
            delete this.activeConnections[id]
        }

        xhr.ontimeout = (error) => {
            reject(error)
            delete this.activeConnections[id]
        }

        let data = new FormData()

        data.append('file_id', this.fileId)
        data.append('chunk_id', id)
        data.append('file', file)

        xhr.send(data)
    })
}

Uploader.prototype.on = function (method, callback) {
    if (typeof callback !== 'function') {
        callback = noop
    }

    this[method] = callback
}

Uploader.prototype.abort = function () {
    Object.keys(this.activeConnections).forEach((id) => {
        this.activeConnections[id].abort()
    })

    this.aborted = true
}

Uploader.prototype.getResult = function () {
    uploadPath({ file_id: this.fileId })
        .then((res) => this.end(null, res))
        .catch((error) => this.complete(error))
}

export const uploader = function () {
    const multiThreadedUploader = new Uploader()

    return {
        options: function (options) {
            multiThreadedUploader.setOptions(options)

            return this
        },

        send: function (file) {
            multiThreadedUploader.setupFile(file)

            return this
        },

        continue: function () {
            multiThreadedUploader.sendNext()
        },

        onProgress: function (callback) {
            multiThreadedUploader.on('onProgress', callback || noop)

            return this
        },

        end: function (callback) {
            multiThreadedUploader.on('end', callback || noop)
            multiThreadedUploader.start()

            return this
        },

        abort: function () {
            multiThreadedUploader.abort()
        },
    }
}

export default uploader
