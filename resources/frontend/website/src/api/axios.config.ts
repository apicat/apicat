import Ajax, { NoMessageAjax } from './Ajax'
import { HTTP_STATUS } from '@/common/constant'
import { ElMessage } from 'element-plus'
import { AxiosRequestConfig } from 'axios'
import { AxiosCancelTokenManager } from '@ac/axios'
import { useUserStore } from '@/stores/user'
import { usePreviewStore } from '@/stores/preview'

export const axiosCancelTokenManager = new AxiosCancelTokenManager()

const initAxiosFilter = () => {
    // 请求拦截器
    Ajax.interceptors.request.use((config: AxiosRequestConfig) => {
        axiosCancelTokenManager.add(config)

        const authorization = useUserStore().token

        if (authorization) {
            config.headers = { ...config.headers, Authorization: authorization }
        }

        return config
    })

    // 响应拦截器
    Ajax.interceptors.response.use((res) => {
        axiosCancelTokenManager.remove(res.config)

        let errorMessage = ''

        !res.data && (res.data = {})

        const status = res.data.status

        // 保存 token信息
        if (res.headers['authorization']) {
            useUserStore().updateToken(res.headers['authorization'])
        }

        switch (status) {
            // 未登录
            case HTTP_STATUS.NO_LOGIN:
                errorMessage = '登录状态已失效，请重新登录'
                axiosCancelTokenManager.removeAll()
                useUserStore().logout()
                break

            // 预览秘钥失效
            case HTTP_STATUS.INVALID_PREVIEW_SECRET:
                errorMessage = res.data.msg || '访问秘钥已失效。'
                axiosCancelTokenManager.removeAll()
                usePreviewStore().goVerificationPage()
                break

            // 正常状态
            case HTTP_STATUS.OK:
            case HTTP_STATUS.NO_PARENT_DIR:
            case HTTP_STATUS.NOT_FOUND:
                break

            // 其他状态
            default:
                errorMessage = res.data.msg || '数据异常，请稍后再试！'
                break
        }

        if (errorMessage) {
            return Promise.reject({ message: errorMessage })
        }

        return res.data
    })

    // 错误拦截
    Ajax.interceptors.response.use(undefined, (error) => {
        let errMessage = ''

        const { code, message } = error || {}
        const err: string = error?.toString?.() ?? ''

        errMessage = message

        if (code === 'ECONNABORTED' && errMessage.indexOf('timeout') !== -1) {
            errMessage = '请求超时，请重试'
        }

        if (err?.includes('Network Error')) {
            errMessage = '网络异常，请重试'
        }

        ElMessage.closeAll()
        errMessage && ElMessage.error(errMessage)

        return Promise.reject(error)
    })

    // NoMessageAjax 请求拦截器
    NoMessageAjax.interceptors.request.use((config) => {
        const authorization = useUserStore().token

        if (authorization) {
            config.headers = { ...config.headers, Authorization: authorization }
        }

        return config
    })

    // 响应拦截器
    NoMessageAjax.interceptors.response.use(
        (res) => {
            let errorMessage

            !res.data && (res.data = {})
            switch (res.data.status) {
                case 0:
                    break
                default:
                    errorMessage = res.data.msg || '数据异常，请稍后再试！'
                    break
            }

            if (errorMessage) {
                return Promise.reject(errorMessage)
            }

            return res.data
        },
        (error) => {
            error.message.indexOf('timeout') !== -1 && ElMessage.error('请求超时，请重试！')
            return Promise.reject(error)
        }
    )
}

export default initAxiosFilter()
