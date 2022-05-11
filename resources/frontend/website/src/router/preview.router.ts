import { compile } from 'path-to-regexp'
import PreviewLayout from '@/layout/PreviewLayout.vue'
import { PREVIEW_DOCUMENT, PREVIEW_DOCUMENT_SECRET, PREVIEW_PROJECT, PREVIEW_PROJECT_SECRET } from './constant'

const ProjectPreview = () => import('../views/preview/ProjectPreview.vue')
const ProjectVerification = () => import('../views/preview/ProjectVerification.vue')

const DocumentPreview = () => import('../views/preview/DocumentPreview.vue')
const DocumentVerification = () => import('../views/preview/DocumentVerification.vue')

const TrashDocumentPreview = () => import('../views/preview/TrashDocumentPreview.vue')

export const toPreviewProjectPath = compile('/app/:project_id')
export const toPreviewTrashDocumentPath = compile('/doc/:project_id/trash_api_preview/:doc_id')

export const documentPreviewRouters = {
    path: '/doc',
    name: PREVIEW_DOCUMENT,
    component: PreviewLayout,
    meta: { ignoreAuth: true },
    children: [
        // 文档 - 密钥校验
        {
            path: '/doc/:doc_id/verification',
            name: PREVIEW_DOCUMENT_SECRET,
            component: DocumentVerification,
            meta: { ignoreAuth: true },
        },
        // 单篇文档预览
        {
            path: '/doc/:doc_id',
            name: `${PREVIEW_DOCUMENT}.detail`,
            component: DocumentPreview,
            meta: { ignoreAuth: true },
        },
    ],
}

export const projectPreviewRouters = {
    path: '/app',
    name: PREVIEW_PROJECT,
    component: PreviewLayout,
    meta: { ignoreAuth: true },
    children: [
        // 项目 - 密钥校验
        {
            path: '/app/:project_id/verification',
            name: PREVIEW_PROJECT_SECRET,
            component: ProjectVerification,
            meta: { ignoreAuth: true },
        },
        // 项目预览 文档路径
        {
            path: '/app/:project_id/:node_id?',
            name: `${PREVIEW_PROJECT}.document`,
            component: ProjectPreview,
            meta: { ignoreAuth: true },
        },
    ],
}

export const trashPreviewRouters = {
    path: '/doc/:project_id/trash_api_preview/:doc_id',
    name: 'preview.trash',
    component: PreviewLayout,
    children: [
        {
            path: '',
            name: 'preview.trash.document',
            component: TrashDocumentPreview,
        },
    ],
}
