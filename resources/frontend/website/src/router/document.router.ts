import { compile } from 'path-to-regexp'
import DocumentLayout from '@/layout/DocumentLayout.vue'
import { DOCUMENT_ROUTE_NAME, DOCUMENT_EDIT_NAME } from './constant'

const DocumentDetail = () => import('@/views/document/DocumentDetail.vue')
const DocumentEditor = () => import('@/views/document/DocumentEditor.vue')
const DECUMENT_DETAIL_PATH = '/editor/:project_id/doc/:node_id?'

export const toDocumentDetailPath = compile(DECUMENT_DETAIL_PATH)

export default {
    path: '/editor',
    name: DOCUMENT_ROUTE_NAME,
    component: DocumentLayout,
    redirect: { name: 'document.api.detail' },
    children: [
        {
            path: DECUMENT_DETAIL_PATH,
            name: 'document.api.detail',
            component: DocumentDetail,
        },
        {
            path: '/editor/:project_id/doc/:node_id/edit',
            name: DOCUMENT_EDIT_NAME,
            component: DocumentEditor,
        },
    ],
}
