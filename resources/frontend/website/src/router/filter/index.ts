import type { Router } from 'vue-router'

import initPermissionsFilter from './permission.filter'
import initDocumentEditFilter from './document.edit.filter'
import initPreviewFilter from './preview.filter'
import initProjectNavFilter from './project.filter'

export const setupRouterFilter = (router: Router) => {
    initPermissionsFilter(router)
    initDocumentEditFilter(router)
    initProjectNavFilter(router)
    initPreviewFilter(router)
}
