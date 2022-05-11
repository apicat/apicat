import { defineStore } from 'pinia'

export const useAppStore = defineStore({
    id: 'app.store',
    state: (): { isShowLoading: boolean } => ({
        isShowLoading: false,
    }),

    actions: {
        showLoading() {
            this.isShowLoading = true
        },

        hideLoading() {
            this.isShowLoading = false
        },
    },
})
