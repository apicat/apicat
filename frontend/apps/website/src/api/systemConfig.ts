import { API_URL, REQUEST_TIMEOUT } from '@/commons'
import axios from 'axios'
import { ElMessage } from 'element-plus'

axios.defaults.timeout = REQUEST_TIMEOUT

const QuietAjax = axios.create({
  baseURL: API_URL,
  headers: {
    Accept: 'application/json, text/plain, */*',
  },
})

QuietAjax.interceptors.response.use((response) => response.data, (error) => {
  const { response = { data: {} } } = error
  const { message } = response.data
  message && ElMessage.error(message)
})


export const saveDBConfig = async (dbConfig: Record<string, any>) => await QuietAjax.put('/config/db', dbConfig)
