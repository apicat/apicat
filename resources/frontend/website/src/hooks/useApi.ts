import { ref } from 'vue'
import { ElMessage } from 'element-plus'

interface UseApiOptions {
    isCatch?: boolean
    isShowMessage?: boolean
    msg?: string
}

export function useApi(execute: any, options?: UseApiOptions): any {
    options = { isCatch: true, isShowMessage: true, ...options }

    const isLoading = ref(false)
    const { isCatch, isShowMessage, msg } = options

    async function call(data: unknown) {
        isLoading.value = true
        try {
            const res = (await execute(data)) || {}
            if (isShowMessage && (msg || res.msg)) {
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
