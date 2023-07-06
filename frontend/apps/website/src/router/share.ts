import { RouteRecordRaw } from 'vue-router'
import PreviewLayout from '@/layouts/PreviewLayout.vue'
import { compile } from 'path-to-regexp'
import { getDocumentDetailPath, getProjectDetailPath } from './project.detail'
import { DOCUMENT_SHARE_PATH, DOCUMENT_SHARE_VALIDATION_NAME, DOCUMENT_SHARE_VALIDATION_PATH, PROJECT_SHARE_VALIDATION_NAME, PROJECT_SHARE_VALIDATION_PATH } from './constant'

const ProjectVerification = () => import('@/views/share/ProjectVerification.vue')
const DocumentVerification = () => import('@/views/share/DocumentVerification.vue')
const DocumentPreview = () => import('@/views/share/DocumentPreview.vue')

// 项目密钥校验
const projectVerificationRoute: RouteRecordRaw = {
  name: PROJECT_SHARE_VALIDATION_NAME,
  path: PROJECT_SHARE_VALIDATION_PATH,
  meta: { ignoreAuth: true },
  component: ProjectVerification,
}

// 文档密钥校验

const documentVerificationRoute: RouteRecordRaw = {
  name: DOCUMENT_SHARE_VALIDATION_NAME,
  path: DOCUMENT_SHARE_VALIDATION_PATH,
  meta: { ignoreAuth: true },
  component: DocumentVerification,
}

// 文档分享详情
const documentShareDetailRoute: RouteRecordRaw = {
  name: 'share.document',
  path: DOCUMENT_SHARE_PATH,
  meta: { ignoreAuth: true },
  component: PreviewLayout,
  children: [
    {
      name: 'share.document.detail',
      path: '',
      component: DocumentPreview,
    },
  ],
}

/**
 * export all routes
 */
export const shareRoutes = [projectVerificationRoute, documentVerificationRoute, documentShareDetailRoute]

// 获取文档密钥校验地址
export const getDocumentVerificationPath = (doc_public_id: string) => compile(DOCUMENT_SHARE_VALIDATION_PATH)({ doc_public_id })
// 获取分享文档详情地址
export const getDocumentShareDetailPath = (doc_public_id: string) => compile(DOCUMENT_SHARE_PATH)({ doc_public_id })
// 获取私有文档分享链接
export const getDocumentPrivateShareLink = (doc_public_id: string) => window.origin + (doc_public_id ? getDocumentShareDetailPath(doc_public_id) : '')
// 获取公共文档分享链接
export const getDocumentPublicShareLink = (project_id: string, doc_id: string) => window.origin + (doc_id ? getDocumentDetailPath(project_id, doc_id) : '')
// 获取项目分享链接
export const getProjectShareLink = (project_public_id: string) => window.origin + getProjectDetailPath(project_public_id)
