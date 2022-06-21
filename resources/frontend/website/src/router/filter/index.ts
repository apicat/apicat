import type { Router } from 'vue-router'

import initPermissionsFilter from './permission.filter'
import initDocumentFilter, { initProjectDetail } from './document.filter'
import initDocumentEditFilter from './document.edit.filter'
import initDocumentDetailFilter from './document.detail.filter'
import initProjectSettingFilter from './project.setting.filter'
// import initPreviewFilter from './preview.filter'

export const setupRouterFilter = (router: Router) => {
    initPermissionsFilter(router)
    initDocumentFilter(router)
    initDocumentEditFilter(router)
    initDocumentDetailFilter(router)
    initProjectSettingFilter(router)

    initProjectDetail(router)
    // initPreviewFilter(router)
}
