import { RouteRecordRaw } from 'vue-router'
import MainLayout from '@/layouts/MainLayout.vue'
import { compile } from 'path-to-regexp'

/**
 * 项目密钥校验
 */
export const projectVerificationRoute: RouteRecordRaw = {
  name: 'proejct.verification.route',
  path: '/project/:project_id/verification',
  meta: { ignoreAuth: true },
  component: MainLayout,
}

/**
 * 文档密钥校验
 */
export const collectionVerificationRoute: RouteRecordRaw = {
  name: 'collection.verification.route',
  path: '/share/:doc_public_id/verification',
  meta: { ignoreAuth: true },
  component: MainLayout,
}

/**
 * 文档密钥校验
 */
const DOCUMENT_SHARE_PATH = '/share/:doc_public_id'
export const shareDocumentRoute: RouteRecordRaw = {
  name: 'collection.detail.share',
  path: DOCUMENT_SHARE_PATH,
  meta: { ignoreAuth: true },
  component: MainLayout,
}

export const getDocumentSharePath = (doc_public_id: string) => compile(DOCUMENT_SHARE_PATH)({ doc_public_id })
