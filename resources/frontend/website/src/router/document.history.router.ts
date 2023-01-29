import DocumentHistoryLayout from '@/layout/DocumentHistoryLayout.vue'
import { DOCUMENT_HISTORY_PATH, DOCUMENT_HISTORY_NAME, DOCUMENT_HISTORY_DETAIL_PATH, DOCUMENT_HISTORY_DETAIL_NAME } from './constant'

const DocumentHistoryDetail = () => import('@/views/document/DocumentHistoryDetail.vue')

export default {
    path: DOCUMENT_HISTORY_PATH,
    name: DOCUMENT_HISTORY_NAME,
    component: DocumentHistoryLayout,
    redirect: { name: DOCUMENT_HISTORY_DETAIL_NAME },
    children: [
        {
            path: DOCUMENT_HISTORY_DETAIL_PATH,
            name: DOCUMENT_HISTORY_DETAIL_NAME,
            meta: { autoHideLoadingLayer: false },
            component: DocumentHistoryDetail,
        },
    ],
}
