import DocumentLayout from '@/layout/DocumentLayout.vue'
import {
    ITERATE_DOCUMENT_DETAIL_NAME,
    ITERATE_DOCUMENT_DETAIL_PATH,
    ITERATE_DOCUMENT_EDIT_PATH,
    ITERATE_DOCUMENT_EDIT_NAME,
    ITERATE_ROUTE_PATH,
    ITERATE_ROUTE_NAME,
} from './constant'

const DocumentDetail = () => import('@/views/document/DocumentDetail.vue')
const DocumentEditor = () => import('@/views/document/DocumentEditor.vue')

export default {
    path: ITERATE_ROUTE_PATH,
    name: ITERATE_ROUTE_NAME,
    component: DocumentLayout,
    redirect: { name: ITERATE_DOCUMENT_DETAIL_NAME },
    children: [
        {
            path: ITERATE_DOCUMENT_DETAIL_PATH,
            name: ITERATE_DOCUMENT_DETAIL_NAME,
            component: DocumentDetail,
        },
        {
            path: ITERATE_DOCUMENT_EDIT_PATH,
            name: ITERATE_DOCUMENT_EDIT_NAME,
            component: DocumentEditor,
        },
    ],
}
