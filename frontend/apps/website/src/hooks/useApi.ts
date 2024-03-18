import { ref } from 'vue'
import { ElMessage } from 'element-plus'
import type { Ref } from 'vue'
import { delay } from '@apicat/shared'

interface UseApiOptions {
  defaultLoadingStatus?: boolean
  isCatch?: boolean
  isShowMessage?: boolean
  msg?: string
}

// 针对异步方法
export function useAsync<T extends any[], R>(
  method: (...args: T) => Promise<R>,
): [Ref<boolean>, (...args: T) => Promise<R | object>] {
  const isLoading = ref(false)
  async function call(...args: T): Promise<R | object> {
    isLoading.value = true
    try {
      const res = await method(...args)
      isLoading.value = false
      return res
    } catch (e) {
      isLoading.value = false
      throw e
    }
  }
  return [isLoading, call]
}

// type ResponseData<R> = ResponseAPI.Response<R> | R

// 定义包裹方法的函数包装器
export function useApi<T extends any[], R>(
  method: (...args: T) => Promise<R>,
  options?: UseApiOptions,
): [Ref<boolean>, (...args: T) => Promise<R | undefined>, Ref<unknown>] {
  const { isCatch = false, isShowMessage = true, msg, defaultLoadingStatus = false } = options || {}

  const isLoading = ref(defaultLoadingStatus || false)
  const isError = ref<unknown>()
  async function call(...args: T): Promise<R | undefined> {
    isLoading.value = true
    try {
      const res = await method(...args)
      isError.value = undefined
      if (res) {
        if (isShowMessage && (msg || (res as any).msg)) {
          ElMessage.closeAll()
          ElMessage.success(msg || (res as any).msg)
        }
      }
      return res
    } catch (e) {
      isError.value = e
      if (!isCatch) throw e
      else console.error(e)
    } finally {
      setTimeout(() => (isLoading.value = false), 0)
    }
  }

  return [isLoading, call, isError]
}

export default useApi
