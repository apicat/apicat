import type { AxiosError, AxiosResponse, InternalAxiosRequestConfig } from 'axios'
import axios from 'axios'
import { ElMessage } from 'element-plus'
import AxiosWrapper from './request'
import { clearShareToken } from './shareToken'
import { useUserStoreWithOut } from '@/store/user'
import { API_URL, REQUEST_TIMEOUT } from '@/commons/constant'
import Storage from '@/commons/storage'

axios.defaults.timeout = REQUEST_TIMEOUT

const baseConfig = {
  baseURL: API_URL,
  headers: {
    Accept: 'application/json, text/plain, */*',
  },
}

export const DefaultAjax = AxiosWrapper.create(baseConfig, {
  onShowErrorMsg(errorMsg) {
    ElMessage.closeAll()
    ElMessage.error(errorMsg)
  },
  onShowSuccessMsg(successMsg) {
    ElMessage.closeAll()
    ElMessage.success(successMsg)
  },
})

export const QuietAjax = axios.create(baseConfig)
export const MockAjax = axios.create({
  ...baseConfig,
  validateStatus() {
    return true
  },
})

function onRequest(config: InternalAxiosRequestConfig): InternalAxiosRequestConfig {
  const token = Storage.get(Storage.KEYS.TOKEN)
  config.headers.Authorization = token ? `Bearer ${token}` : ''
  return config
}

function onErrorResponse(error: AxiosError | Error): Promise<AxiosError> {
  if (axios.isAxiosError(error)) {
    const useUserStore = useUserStoreWithOut()

    const { response = { data: {} } } = error
    const { message } = response.data
    const { status } = (error.response as AxiosResponse) ?? {}

    let errorMsg = ''
    switch (status) {
      case 401:
        // 其他方法
        switch (response.data.action) {
          case 'login':
            errorMsg = message || 'Token Expired'
            useUserStore.logout()
            break

          case 'verify':
            errorMsg = message || 'Share Token Expired'
            clearShareToken()
            break
        }
        break

      case 429:
      // case 500:
        errorMsg = message || 'Server Error'
        break
    }

    ElMessage.closeAll()
    errorMsg && !axios.isCancel(error) && ElMessage.error(errorMsg)
  }

  return Promise.reject(error)
}

DefaultAjax.addRequestInterceptor(onRequest, onErrorResponse)
DefaultAjax.addResponseInterceptor((response: AxiosResponse) => response.data, onErrorResponse)

QuietAjax.interceptors.request.use(onRequest, onErrorResponse)
QuietAjax.interceptors.response.use((response: AxiosResponse) => response.data, onErrorResponse)

MockAjax.interceptors.request.use(onRequest, onErrorResponse)

// 默认请求实例
export default DefaultAjax

// blob专用
export const RawAjax = axios.create()

RawAjax.interceptors.request.use(onRequest)
RawAjax.interceptors.response.use(
  async (response) => {
    return response.data
  },
  (error) => {
    const errMsg = error.response.status + error?.response?.data?.msg
    ElMessage.closeAll()
    ElMessage.error(errMsg)
    error.message = errMsg
    return Promise.reject(error)
  },
)
