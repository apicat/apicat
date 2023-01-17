import type { Router } from 'vue-router'

import initPermissionsFilter from './permission.filter'
import initProjectNavFilter from './project.setting.filter'

import initDocumentFilter, { initProjectDetailFilter } from './document.filter'
import initDocumentEditFilter from './document.edit.filter'
import initDocumentDetailFilter from './document.detail.filter'

import initIterateFilter from './iterate.filter'
import initIterateDetailFilter from './iterate.detail.filter'
import initIterateEditFilter from './iterate.edit.filter'

import initPreviewFilter from './preview.filter'

export const setupRouterFilter = (router: Router) => {
    initPermissionsFilter(router)

    initDocumentFilter(router)
    initDocumentEditFilter(router)
    initDocumentDetailFilter(router)

    initProjectDetailFilter(router)
    initPreviewFilter(router)
    initProjectNavFilter(router)

    initIterateFilter(router)
    initIterateDetailFilter(router)
    initIterateEditFilter(router)
}
