import { compile } from 'path-to-regexp'
import DocumentLayout from '@/layout/DocumentLayout.vue'
import { DOCUMENT_ROUTE_NAME, DOCUMENT_EDIT_NAME, DOCUMENT_DETAIL_NAME, PROJECT_PREVIEW_PATH, DECUMENT_DETAIL_PATH, DECUMENT_EDIT_PATH } from './constant'

const DocumentDetail = () => import('@/views/document/DocumentDetail.vue')
const DocumentEditor = () => import('@/views/document/DocumentEditor.vue')

export const toDocumentDetailPath = compile(DECUMENT_DETAIL_PATH)
export const toDocumentEditPath = compile(DECUMENT_EDIT_PATH)

export default {
    path: PROJECT_PREVIEW_PATH,
    name: DOCUMENT_ROUTE_NAME,
    component: DocumentLayout,
    redirect: { name: DOCUMENT_DETAIL_NAME },
    children: [
        {
            path: DECUMENT_DETAIL_PATH,
            name: DOCUMENT_DETAIL_NAME,
            meta: { ignoreAuth: true },
            component: DocumentDetail,
        },
        {
            path: DECUMENT_EDIT_PATH,
            name: DOCUMENT_EDIT_NAME,
            component: DocumentEditor,
        },
    ],
}
