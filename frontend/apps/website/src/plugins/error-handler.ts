import type { App } from 'vue'
import { useGlobalLoading } from '@/hooks/useGlobalLoading'

export default (app: App) => {
  app.config.errorHandler = (err, vm, info) => {
    // todo send error info to server
    console.error(err, vm, info)
    useGlobalLoading().hideGlobalLoading()
  }
}
