import { useShareStore } from '@/store/share'
import { useUserStoreWithOut } from '@/store/user'
import axios, { AxiosError, AxiosResponse, InternalAxiosRequestConfig } from 'axios'
import { API_URL, PERMISSION_CHANGE_CODE, REQUEST_TIMEOUT, RESPONSE_NOT_FOUND_MAPS, RESPONSE_UNAUTHORIZED_MAPS } from '@/commons/constant'
import { ElMessage, ElMessageBox } from 'element-plus'
import Storage from '@/commons/storage'
import { LOGIN_PATH, NOT_FOUND_PATH, NO_PERMISSION_PATH, router } from '@/router'
import { i18n } from '@/i18n'
import { ShareSecretKeyError, TargetMemberPermissionError } from './error'
import useProjectStore from '@/store/project'

axios.defaults.timeout = REQUEST_TIMEOUT

const baseConfig = {
  baseURL: API_URL,
  headers: {
    Accept: 'application/json, text/plain, */*',
  },
}

let isShowPermissionChangeModal = false

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
  const shareStore = useShareStore()
  const projectStore = useProjectStore()

  let errorMsg = ''
  if (axios.isAxiosError(error)) {
    const { response = { data: {} } } = error
    const { code, message } = response.data
    const { status } = (error.response as AxiosResponse) ?? {}
    const { t } = i18n.global as any
    errorMsg = message

    switch (status) {
      case 401: // 未登录
        if (code === RESPONSE_UNAUTHORIZED_MAPS.LOGIN_TOKEN_EXPIRED_OR_ERROR) {
          useUserStore.logout()
          setTimeout(() => router.replace(LOGIN_PATH), 0)
          return Promise.reject(error)
        }

        // 分享文档｜项目密钥失效
        if (code === RESPONSE_UNAUTHORIZED_MAPS.PROJECT_OR_DOCUMENT_SECRET_TOKEN_EXPIRED_OR_ERROR) {
          errorMsg = response.data.message
          error = new ShareSecretKeyError()
        }

        break

      case 403: // 无权限
        errorMsg = ''

        if (isShowPermissionChangeModal) {
          ElMessage.closeAll()
          return Promise.reject(error)
        }

        code === PERMISSION_CHANGE_CODE.USER_PREMISSION_ERROR &&
          ElMessageBox({
            type: 'warning',
            title: t('app.tips.permissionChangeTitle'),
            message: t('app.user.tips.permissionChange'),
            'show-close': false,
            'close-on-click-modal': false,
            'close-on-press-escape': false,
            'show-cancel-button': false,
            confirmButtonText: t('app.common.confirm'),
            callback() {
              isShowPermissionChangeModal = true
              location.reload()
            },
          } as any)

        code === PERMISSION_CHANGE_CODE.MEMBER_PREMISSION_ERROR &&
          ElMessageBox({
            type: 'warning',
            title: t('app.tips.permissionChangeTitle'),
            message: t('app.project.tips.permissionChange'),
            'show-close': false,
            'close-on-click-modal': false,
            'close-on-press-escape': false,
            'show-cancel-button': false,
            confirmButtonText: t('app.common.confirm'),
            callback() {
              isShowPermissionChangeModal = true
              location.reload()
            },
          } as any)

        if (code === PERMISSION_CHANGE_CODE.TARGET_MEMBER_PREMISSION_ERROR) {
          error = new TargetMemberPermissionError()
        }

        // 重定向403页面
        if (code === PERMISSION_CHANGE_CODE.REDIRECT_UNAUTHORIZED_PAGE) {
          setTimeout(() => router.replace(NO_PERMISSION_PATH), 0)
        }
        break

      case 400: // bad request
      case 404: // not found
        if (code === RESPONSE_NOT_FOUND_MAPS.REDIRECT_NOT_FOUND_PAGE) {
          setTimeout(() => router.replace(NOT_FOUND_PATH), 0)
        }
        break

      default:
        errorMsg = errorMsg || error.message || 'server error'
        break
    }
  }

  ElMessage.closeAll()
  errorMsg && ElMessage.error(errorMsg)
  return Promise.reject(error)
}

DefaultAjax.interceptors.request.use(onRequest, onErrorResponse)
DefaultAjax.interceptors.response.use((response: AxiosResponse) => {
  if (response.status > 200) {
    ElMessage.closeAll()
    ElMessage.success(response.data.message || 'success')
  }
  return response.data
}, onErrorResponse)

QuietAjax.interceptors.request.use(onRequest, onErrorResponse)
QuietAjax.interceptors.response.use((response: AxiosResponse) => response.data, onErrorResponse)

MockAjax.interceptors.request.use(onRequest, onErrorResponse)

// 默认请求实例
export default DefaultAjax
