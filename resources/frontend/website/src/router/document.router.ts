import { compile } from 'path-to-regexp'
import DocumentLayout from '@/layout/DocumentLayout.vue'
import { DOCUMENT_ROUTE_NAME, DOCUMENT_EDIT_NAME, DOCUMENT_DETAIL_NAME, PROJECT_DETAIL_PATH, DOCUMENT_DETAIL_PATH, DOCUMENT_EDIT_PATH } from './constant'

const DocumentDetail = () => import('@/views/document/DocumentDetail.vue')
const DocumentEditor = () => import('@/views/document/DocumentEditor.vue')

export const toDocumentDetailPath = compile(DOCUMENT_DETAIL_PATH)
export const toDocumentEditPath = compile(DOCUMENT_EDIT_PATH)

export default {
    path: PROJECT_DETAIL_PATH,
    name: DOCUMENT_ROUTE_NAME,
    component: DocumentLayout,
    redirect: { name: DOCUMENT_DETAIL_NAME },
    children: [
        {
            path: DOCUMENT_DETAIL_PATH,
            name: DOCUMENT_DETAIL_NAME,
            meta: { ignoreAuth: true },
            component: DocumentDetail,
        },
        {
            path: DOCUMENT_EDIT_PATH,
            name: DOCUMENT_EDIT_NAME,
            component: DocumentEditor,
        },
    ],
}
