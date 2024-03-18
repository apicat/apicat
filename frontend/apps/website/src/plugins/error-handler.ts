import type { App } from 'vue'

export default (app: App) => {
  app.config.errorHandler = (err, vm, info) => {
    // todo send error info to server
    console.error(err, vm, info)
  }
}
