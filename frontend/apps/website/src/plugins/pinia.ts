import { createPinia } from 'pinia'
import type { Router } from 'vue-router'

export const pinia = createPinia()

export const setupPiniaWithRouter = (router: Router) => {
  pinia.use(({ store }) => {
    store.$router = markRaw(router)
  })
}

export default pinia
