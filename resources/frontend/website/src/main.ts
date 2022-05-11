import { createApp } from 'vue'
import { createPinia } from 'pinia'

import App from './App.vue'
import router, { setupRouterFilter } from './router'

import elementSetup from './element.setup'
import ClipboardHelper from './components/ClipboardHelper'

import './assets/stylesheet/index.scss'
import './api/axios.config'

const app = createApp(App)

app.use(createPinia())

app.use(router)
app.use(ClipboardHelper)

setupRouterFilter(router)

elementSetup(app)

app.mount('#app')
