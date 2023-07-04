import { RouteRecordRaw } from 'vue-router'
import PreviewLayout from '@/layouts/PreviewLayout.vue'
import { compile } from 'path-to-regexp'
import { getDocumentDetailPath, getProjectDetailPath } from './project.detail'

const ProjectVerification = () => import('@/views/share/ProjectVerification.vue')
const DocumentVerification = () => import('@/views/share/DocumentVerification.vue')
const DocumentPreview = () => import('@/views/share/DocumentPreview.vue')

// 项目密钥校验
const projectVerificationRoute: RouteRecordRaw = {
  name: 'share.proejct.verification',
  path: '/project/:project_id/verification',
  meta: { ignoreAuth: true },
  component: ProjectVerification,
}

// 文档密钥校验
const documentVerificationRoute: RouteRecordRaw = {
  name: 'share.document.verification',
  path: '/share/:doc_public_id/verification',
  meta: { ignoreAuth: true },
  component: DocumentVerification,
}

// 文档分享详情
const DOCUMENT_SHARE_PATH = '/share/:doc_public_id'
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

// 获取私有文档分享链接
export const getDocumentPrivateShareLink = (doc_public_id: string) => window.origin + (doc_public_id ? compile(DOCUMENT_SHARE_PATH)({ doc_public_id }) : '')
// 获取公共文档分享链接
export const getDocumentPublicShareLink = (project_id: string, doc_id: string) => window.origin + (doc_id ? getDocumentDetailPath(project_id, doc_id) : '')
// 获取项目分享链接
export const getProjectShareLink = (project_public_id: string) => window.origin + getProjectDetailPath(project_public_id)
