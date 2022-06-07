import { App } from 'vue'

import { ElMessage, ElMessageBox, ElLoading } from 'element-plus'
import AsyncMessageBox from '@/components/AsyncMessageBox'

import 'element-plus/theme-chalk/base.css'
import 'element-plus/theme-chalk/el-message.css'
import 'element-plus/theme-chalk/el-message-box.css'
import 'element-plus/theme-chalk/el-loading.css'
import 'element-plus/theme-chalk/el-tree.css'

/**
 * 处理element-ui 无法按需导入的组件
 */
export default (app: App) => {
    app.use(ElMessage)
    app.use(ElMessageBox)
    app.use(AsyncMessageBox)
    app.use(ElLoading)

    // alias for iview $Message
    app.config.globalProperties.$Message = ElMessage
}
