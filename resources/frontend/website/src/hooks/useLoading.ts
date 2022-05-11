// import { useAppStore } from '../stores/app'
// export const showLoading = () => useAppStore().showLoading()
// export const hideLoading = () => useAppStore().hideLoading()

import { ElLoading } from 'element-plus'

let globalLoading: any = null

export const showLoading = () => {
    globalLoading = ElLoading.service({
        lock: true,
        background: 'rgba(255, 255, 255, 1)',
    })
}

export const hideLoading = () => {
    globalLoading && globalLoading.close()
}
