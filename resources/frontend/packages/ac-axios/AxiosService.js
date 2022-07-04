import axios from 'axios'
import { isFunction } from '@natosoft/shared'
import AxiosPendingManager from './AxiosCancelTokenManager'

export default class AxiosService {
    constructor(axiosConfig = {}, options = {}) {
        this.options = { ...this.defaultOptions, ...options }

        this.axios = axios.create(axiosConfig)

        this.axiosPendingManager = new AxiosPendingManager()

        // 设置拦截器
        this.setupInterceptors()
    }

    get defaultOptions() {
        return {
            // 请求
            requestInterceptor: null,
            responseInterceptor: null,
            errorHandler: null,
            // 是否取消重复请求
            isCancelRepeatRequest: true,
        }
    }

    setupInterceptors() {
        const { requestInterceptor, responseInterceptor, errorHandler, isCancelRepeatRequest } = this.options

        this.axios.interceptors.request.use((config) => {
            isCancelRepeatRequest && this.axiosPendingManager.add(config)
            if (requestInterceptor && isFunction(requestInterceptor)) {
                requestInterceptor(config)
            }
            return config
        })

        this.axios.interceptors.response.use((res) => {
            isCancelRepeatRequest && this.axiosPendingManager.remove()
            if (responseInterceptor && isFunction(responseInterceptor)) {
                res = responseInterceptor(res)
            }
            return res
        })

        errorHandler && isFunction(errorHandler) && this.axios.interceptors.response.use(undefined, errorHandler)
    }

    get(url, axiosRequestConfig) {
        return this.request({ ...axiosRequestConfig, url, method: 'get' })
    }

    post(url, axiosRequestConfig) {
        return this.request({ ...axiosRequestConfig, url, method: 'post' })
    }

    delete(url, axiosRequestConfig) {
        return this.request({ ...axiosRequestConfig, url, method: 'delete' })
    }

    put(url, axiosRequestConfig) {
        return this.request({ ...axiosRequestConfig, url, method: 'put' })
    }

    request(axiosRequestConfig) {
        return new Promise(function (resolve, reject) {
            this.axios
                .request(axiosRequestConfig)
                .then((res) => resolve(res))
                .catch((error) => reject(error))
        })
    }
}
