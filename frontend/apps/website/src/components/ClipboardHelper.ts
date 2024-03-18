import ClipboardJS from 'clipboard'
import { ElMessage } from 'element-plus'
import { trim } from 'lodash-es'
import { useLocaleStoreWithOut } from '@/store/locale'

const install = function (app: any) {
  app.mixin({
    mounted() {
      const { i18n } = useLocaleStoreWithOut()
      if (app.__clipboard__)
        app.__clipboard__.destroy()

      app.__clipboard__ = new ClipboardJS('.copy_text', {
        text(trigger: HTMLElement) {
          return trim(
            trigger.getAttribute('data-text') || trigger.innerText || '',
          )
        },
      })

      app.__clipboard__.on('success', () => {
        ElMessage.closeAll()
        ElMessage.success((i18n.global as any).t('app.tips.copyed'))
      })
    },
  })
}

export default { install }
