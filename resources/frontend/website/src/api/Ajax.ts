import axios from 'axios'
import { API_URL, REQUEST_TIMEOUT } from '@ac/shared'

const config = {
    baseURL: API_URL,
    headers: {
        Accept: 'application/json, text/plain, */*',
    },
}

axios.defaults.timeout = REQUEST_TIMEOUT

// 默认请求实例
export default axios.create(config)

export const NoMessageAjax = axios.create(config)
