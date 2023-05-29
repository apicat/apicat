import { useUserStoreWithOut } from '@/store/user'
import axios, { AxiosError, AxiosResponse, InternalAxiosRequestConfig } from 'axios'
import { API_URL, REQUEST_TIMEOUT } from '@/commons/constant'
import { ElMessage } from 'element-plus'
import Storage from '@/commons/storage'
import { LOGIN_PATH, router } from '@/router'

axios.defaults.timeout = REQUEST_TIMEOUT

const baseConfig = {
  baseURL: API_URL,
  headers: {
    Accept: 'application/json, text/plain, */*',
  },
}

export const DefaultAjax = axios.create(baseConfig)
export const QuietAjax = axios.create(baseConfig)
export const MockAjax = axios.create({
  ...baseConfig,
  validateStatus: function () {
    return true
  },
})

const onRequest = (config: InternalAxiosRequestConfig): InternalAxiosRequestConfig => {
  const token = Storage.get(Storage.KEYS.TOKEN)
  config.headers.Authorization = token ? `Bearer ${token}` : ''
  return config
}

const onErrorResponse = (error: AxiosError | Error): Promise<AxiosError> => {
  const useUserStore = useUserStoreWithOut()

  if (axios.isAxiosError(error)) {
    const { response = { data: {} } } = error
    const { status } = (error.response as AxiosResponse) ?? {}

    // "Login required"
    if (status === 401) {
      useUserStore.logout()
      setTimeout(() => router.replace(LOGIN_PATH), 0)
    }

    ElMessage.error(response.data.message || 'server error')
  } else {
    ElMessage.error(error.message || 'server error')
  }

  return Promise.reject(error)
}

DefaultAjax.interceptors.request.use(onRequest, onErrorResponse)
DefaultAjax.interceptors.response.use((response: AxiosResponse) => {
  if (response.status > 200) {
    ElMessage.success(response.data.message || 'success')
  }
  return response.data
}, onErrorResponse)

QuietAjax.interceptors.request.use(onRequest, onErrorResponse)
QuietAjax.interceptors.response.use((response: AxiosResponse) => response.data, onErrorResponse)

MockAjax.interceptors.request.use(onRequest, onErrorResponse)
MockAjax.interceptors.response.use((response: AxiosResponse) => response.data)

// 默认请求实例
export default DefaultAjax
