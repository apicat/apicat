import { ElLoading } from 'element-plus'

let globalLoading: any = null

export const showLoading = () => {
    if (globalLoading && globalLoading.visible.value) {
        return
    }
    globalLoading = ElLoading.service({
        lock: true,
        background: 'rgba(255, 255, 255, 1)',
    })
}

export const hideLoading = () => {
    globalLoading && globalLoading.close()
}
