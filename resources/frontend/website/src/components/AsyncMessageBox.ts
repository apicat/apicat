import { ElMessageBox } from 'element-plus'

const useAsyncMsgBox =
    ($msgbox = ElMessageBox) =>
    ({ title, message, content, cb, onOk }: any) => {
        cb = cb || onOk
        message = message || content
        return $msgbox({
            title: title,
            message,
            showCancelButton: true,
            confirmButtonText: '确定',
            cancelButtonText: '取消',
            beforeClose: function (action, instance, done) {
                if (action === 'confirm') {
                    instance.confirmButtonLoading = true
                    Promise.resolve(cb(done))
                        .then(() => done())
                        .catch((e) => e)
                        .finally(() => {
                            instance.confirmButtonLoading = false
                        })
                } else {
                    done()
                }
            },
        })
    }

export const AsyncMsgBox = useAsyncMsgBox()

export default {
    install: function (app: any) {
        const $msgbox = app.config.globalProperties.$msgbox || ElMessageBox
        app.config.globalProperties.$asyncMsgBox = useAsyncMsgBox($msgbox)
    },
}
