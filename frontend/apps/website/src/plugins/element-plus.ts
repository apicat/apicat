import type { App } from 'vue'
import ElementPlus, { ElMessage } from 'element-plus'
import AsyncMessageBox from '@/components/AsyncMessageBox'

export default (app: App) => {
  app.use(ElementPlus)
  app.use(AsyncMessageBox)
  app.config.globalProperties.$Message = ElMessage
}

// import type { App } from 'vue'
// import { ElLoading, ElMessage, ElMessageBox } from 'element-plus'
// import 'element-plus/theme-chalk/src/message-box.scss'
// import 'element-plus/theme-chalk/src/message.scss'
// import 'element-plus/theme-chalk/src/loading.scss'

// import AsyncMessageBox from '@/components/AsyncMessageBox'

// export default (app: App) => {
//   app.use(ElMessage)
//   app.use(ElMessageBox)
//   app.use(AsyncMessageBox)
//   app.use(ElLoading)
//   app.config.globalProperties.$Message = ElMessage
// }
