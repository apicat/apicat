import 'uno.css'
import '@/styles/reset.scss'
import { createApp } from 'vue'
import elementPlus from './plugins/element-plus'
import errorHandler from './plugins/error-handler'

import App from './DBConfigApp.vue'

const run = async () => {
  const app = createApp(App)
  app.use(elementPlus)
  app.use(errorHandler)
  app.mount('#app')
}

run()
