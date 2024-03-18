import '@/styles/element/index.scss'
import '@apicat/components/dist/style.css'
import '@/styles/main.scss'
import 'uno.css'
import '@/assets/iconfont/iconfont'

import { createApp } from 'vue'
import App from './App.vue'

import { elementPlus, errorHandler, pinia, setupPiniaWithRouter } from './plugins'
import clipboardHelper from '@/components/ClipboardHelper'
import limitInput from '@/directives/LimitInput'
import router, { setupRouterFilter } from '@/router'
import { useLocaleStoreWithOut } from '@/store/locale'

async function run() {
  const app = createApp(App)

  const { i18n, initLocale } = useLocaleStoreWithOut()
  await initLocale()

  app.use(i18n)
  app.use(pinia)
  app.use(elementPlus)
  app.use(errorHandler)
  app.use(router)
  app.use(clipboardHelper)
  app.use(limitInput)
  setupRouterFilter(router)
  setupPiniaWithRouter(router)

  app.mount('#app')
}

run()
