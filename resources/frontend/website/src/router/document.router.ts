import { compile } from 'path-to-regexp'
import DocumentLayout from '@/layout/DocumentLayout.vue'
import { DOCUMENT_ROUTE_NAME, DOCUMENT_EDIT_NAME, DOCUMENT_DETAIL_NAME } from './constant'

const DocumentDetail = () => import('@/views/document/DocumentDetail.vue')
const DocumentEditor = () => import('@/views/document/DocumentEditor.vue')
const DECUMENT_DETAIL_PATH = '/editor/:project_id/doc/:node_id?'
const DECUMENT_EDIT_PATH = '/editor/:project_id/doc/:node_id/edit'

export const toDocumentDetailPath = compile(DECUMENT_DETAIL_PATH)
export const toDocumentEditPath = compile(DECUMENT_EDIT_PATH)

export default {
    path: '/editor/:project_id',
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
