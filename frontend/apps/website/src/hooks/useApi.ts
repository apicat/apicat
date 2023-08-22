import { ref } from 'vue'
import { ElMessage } from 'element-plus'

interface UseApiOptions {
  isCatch?: boolean
  isShowMessage?: boolean
  msg?: string
}

export function useApi(execute: any, options?: UseApiOptions): any {
  options = { isCatch: false, isShowMessage: true, ...options }

  const isLoading = ref(false)

  const { isCatch, isShowMessage, msg } = options

  async function call(...args: any[]) {
    isLoading.value = true
    try {
      const _data = isRef(args[0]) ? unref(args[0]) : args[0]
      const res = (await execute(toRaw(_data), ...args.slice(1))) || {}
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

  return [isLoading, call]
}

export default useApi
