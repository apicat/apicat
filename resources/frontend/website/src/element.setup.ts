import { App } from 'vue'

// import { ElMessage, ElMessageBox, ElLoading } from 'element-plus'
// import 'element-plus/theme-chalk/el-message.css'
// import 'element-plus/theme-chalk/el-message-box.css'
// import AsyncMessageBox from '@/components/AsyncMessageBox'

import ElementPlus, { ElMessage } from 'element-plus'
import 'element-plus/dist/index.css'

/**
 * 处理element-ui 无法按需导入的组件
 */
export default (app: App) => {
    // app.use(ElMessage)
    // app.use(ElMessageBox)
    // app.use(AsyncMessageBox)
    // app.use(ElLoading)
    // alias for iview $Message
    app.config.globalProperties.$Message = ElMessage

    app.use(ElementPlus)
}
