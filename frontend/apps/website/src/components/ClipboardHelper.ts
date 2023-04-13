import ClipboardJS from 'clipboard'
import { ElMessage } from 'element-plus'
import { trim } from 'lodash-es'

const install = function (app: any) {
  app.mixin({
    mounted: function () {
      if (app.__clipboard__) {
        app.__clipboard__.destroy()
      }

      app.__clipboard__ = new ClipboardJS('.copy_text', {
        text: function (trigger: HTMLElement) {
          return trim(trigger.getAttribute('data-text') || trigger.innerText || '')
        },
      })

      app.__clipboard__.on('success', () => {
        ElMessage.closeAll()
        ElMessage.success('复制成功')
      })
    },
  })
}

export default { install }
