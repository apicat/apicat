import { ElMessageBox } from 'element-plus'
import type { AppContext } from 'vue'

function createAsyncMsgBox($msgbox = ElMessageBox) {
  return async ({ title, message, content, cb, onOk, ...rest }: any, context?: AppContext) => {
    cb = cb || onOk
    message = message || content
    try {
      return await $msgbox(
        {
          showCancelButton: true,
          ...rest,
          title,
          message,
          beforeClose(action, instance, done) {
            if (action === 'confirm') {
              instance.confirmButtonLoading = true
              Promise.resolve(cb(done))
                .then(() => done())
                .catch(e => e)
                .finally(() => {
                  instance.confirmButtonLoading = false
                })
            }
            else {
              done()
            }
          },
        },
        context,
      )
    }
    catch {}
  }
}

export const AsyncMsgBox = createAsyncMsgBox()

export default {
  install(app: any) {
    const $msgbox = app.config.globalProperties.$msgbox || ElMessageBox
    app.config.globalProperties.$asyncMsgBox = createAsyncMsgBox($msgbox)
  },
}
