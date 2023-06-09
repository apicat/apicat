import { ElMessageBox } from 'element-plus'

const createAsyncMsgBox =
  ($msgbox = ElMessageBox) =>
  ({ title, message, content, cb, onOk, ...rest }: any) => {
    cb = cb || onOk
    message = message || content
    return $msgbox({
      showCancelButton: true,
      ...rest,
      title: title,
      message,
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
    }).catch(() => {
      //
    })
  }

export const AsyncMsgBox = createAsyncMsgBox()

export default {
  install: function (app: any) {
    const $msgbox = app.config.globalProperties.$msgbox || ElMessageBox
    app.config.globalProperties.$asyncMsgBox = createAsyncMsgBox($msgbox)
  },
}
