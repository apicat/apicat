import axios from 'axios'
import type {
  AxiosInstance,
  AxiosRequestConfig,
  AxiosResponse,
  CreateAxiosDefaults,
  InternalAxiosRequestConfig,
} from 'axios'
import { BadRequestError, NoPermissionError, NotFoundError, UnauthorizedError } from './error'
import { REQUEST_TIMEOUT } from '@/commons'

interface AxiosWrapperOptions {
  maxConcurrents?: number
  onShowErrorMsg?: (msg: string) => void
  onShowSuccessMsg?: (msg: string) => void
}

interface RequestOptions {
  isShowErrorMsg?: boolean
  isShowSuccessMsg?: boolean
  successMsg?: string
  failureMsg?: string
}

class AxiosWrapper {
  private instance: AxiosInstance
  private pendingRequests: Map<string, AbortController>
  private onShowErrorMsg?: (msg: string) => void
  private onShowSuccessMsg?: (msg: string) => void

  private constructor(axiosConfig: CreateAxiosDefaults = {}, wrapperOptions?: AxiosWrapperOptions) {
    const { baseURL = '', timeout = REQUEST_TIMEOUT } = axiosConfig
    const { onShowErrorMsg, onShowSuccessMsg } = wrapperOptions || {}

    this.instance = axios.create({ baseURL, timeout })
    this.pendingRequests = new Map()
    this.onShowErrorMsg = onShowErrorMsg
    this.onShowSuccessMsg = onShowSuccessMsg
  }

  public static create(axiosConfig?: CreateAxiosDefaults, wrapperOptions?: AxiosWrapperOptions): AxiosWrapper {
    return new AxiosWrapper(axiosConfig, wrapperOptions)
  }

  public addRequestInterceptor(
    onFulfilled?: (config: InternalAxiosRequestConfig) => InternalAxiosRequestConfig,
    onRejected?: (error: any) => any,
  ): number {
    return this.instance.interceptors.request.use(onFulfilled, onRejected)
  }

  public removeRequestInterceptor(interceptorId: number): void {
    this.instance.interceptors.request.eject(interceptorId)
  }

  public addResponseInterceptor(
    onFulfilled?: (response: AxiosResponse) => AxiosResponse | Promise<AxiosResponse>,
    onRejected?: (error: any) => any,
  ): number {
    return this.instance.interceptors.response.use(onFulfilled, onRejected)
  }

  public removeResponseInterceptor(interceptorId: number): void {
    this.instance.interceptors.response.eject(interceptorId)
  }

  private getRequestKey(config: AxiosRequestConfig): string {
    const { url, params, data } = config
    return `${url}-${JSON.stringify(params)}-${JSON.stringify(data)}`
  }

  public async request<U = any>(config: AxiosRequestConfig, options?: RequestOptions): Promise<U> {
    const requestId = this.getRequestKey(config)
    const { isShowErrorMsg = true, isShowSuccessMsg = false, successMsg = 'Success', failureMsg } = options || {}
    // cancel repeated request
    this.cancelRequest(requestId)

    const controller = new AbortController()
    config.signal = config.signal || controller.signal
    this.pendingRequests.set(requestId, controller)

    try {
      const response = (await this.instance.request(config)) as any
      const message = response?.message || (response as any)?.data?.message || successMsg
      isShowSuccessMsg && this.onShowSuccessMsg && message && this.onShowSuccessMsg(message)
      return response as U
    }
    catch (error: any) {
      // network error
      if (axios.isAxiosError(error) && error?.code === 'ERR_NETWORK')
        this.onShowErrorMsg && this.onShowErrorMsg('Network Error')

      // timeout error
      if (axios.isAxiosError(error) && error?.code === 'ECONNABORTED')
        this.onShowErrorMsg && this.onShowErrorMsg('Request Timeout, Please Try Again Later')

      // http error
      if (axios.isAxiosError(error) && error.response) {
        const { status, data = {} } = error.response || {}
        if (isShowErrorMsg) {
          // show error message
          const errorMsg = data?.message || failureMsg || 'Operation failed.'
          errorMsg && this.onShowErrorMsg && this.onShowErrorMsg(errorMsg)
        }

        // bad request
        if (status === 400)
          throw new BadRequestError(error.response.data)
        else if (status === 403)
          throw new NoPermissionError(error.response.data)
        else if (status === 404)
          throw new NotFoundError(error.response.data)
        else if (status === 401)
          throw new UnauthorizedError(error.response.data)
      }

      throw error
    }
    finally {
      this.pendingRequests.delete(requestId)
    }
  }

  public get<T = any>(url: string, config?: AxiosRequestConfig, options?: RequestOptions) {
    return this.request<T>({ ...config, url, method: 'get' }, options)
  }

  public post<T = any>(url: string, data?: any, options?: RequestOptions, config?: AxiosRequestConfig) {
    return this.request<T>({ ...config, url, method: 'post', data }, options)
  }

  public put<T = any>(url: string, data?: any, options?: RequestOptions, config?: AxiosRequestConfig) {
    return this.request<T>({ ...config, url, method: 'put', data }, options)
  }

  public patch<T = any>(url: string, data?: any, options?: RequestOptions, config?: AxiosRequestConfig) {
    return this.request<T>({ ...config, url, method: 'patch', data }, options)
  }

  public delete<T = any>(url: string, data?: any, options?: RequestOptions, config?: AxiosRequestConfig) {
    return this.request<T>({ ...config, url, method: 'delete', data }, options)
  }

  public cancelRequest(requestId: string): void {
    if (this.pendingRequests.has(requestId)) {
      const controller = this.pendingRequests.get(requestId)
      controller?.abort()
      this.pendingRequests.delete(requestId)
    }
  }

  public cancelAllRequests(): void {
    this.pendingRequests.forEach((controller) => {
      controller?.abort()
    })
    this.pendingRequests.clear()
  }

  public getUri(config: AxiosRequestConfig<any>) {
    return this.instance.getUri(config)
  }
}

export default AxiosWrapper
