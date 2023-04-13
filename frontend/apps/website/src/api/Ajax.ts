import axios, { AxiosResponse } from 'axios'
import { API_URL, REQUEST_TIMEOUT } from '@/commons/constant'
import { ElMessage } from 'element-plus'

axios.defaults.timeout = REQUEST_TIMEOUT

const baseConfig = {
  baseURL: API_URL,
  headers: {
    Accept: 'application/json, text/plain, */*',
  },
}

const handleError = (error: any) => {
  const { response } = error
  ElMessage.error(response.data.message || 'server error')
  return Promise.reject(error)
}

const instance = axios.create(baseConfig)
export const QuietAjax = axios.create(baseConfig)

instance.interceptors.response.use((response: AxiosResponse) => {
  if (response.status > 200) {
    ElMessage.success(response.data.message || 'success')
  }
  return response.data
}, handleError)

QuietAjax.interceptors.response.use((response: AxiosResponse) => response.data, handleError)

// 默认请求实例
export default instance
