import { ref } from 'vue'
import { ElMessage } from 'element-plus'

interface UseApiOptions {
  isCatch?: boolean
  isShowMessage?: boolean
  msg?: string
}

export function useApi(execute: any, options?: UseApiOptions): any {
  options = { isCatch: false, isShowMessage: true, ...options }

  const isLoading = ref(true)

  const { isCatch, isShowMessage, msg } = options

  async function call(data: unknown) {
    isLoading.value = true
    try {
      const _data = isRef(data) ? unref(data) : data
      const res = (await execute(toRaw(_data))) || {}
      if (isShowMessage && (msg || res.msg)) {
        ElMessage.closeAll()
        ElMessage.success(msg || res.msg)
      }
      return res
    } catch (e) {
      if (!isCatch) {
        throw e
      } else {
        return {}
      }
    } finally {
      isLoading.value = false
    }
  }

  return (): [Ref<boolean>, any] => [isLoading, call]
}

export default useApi
