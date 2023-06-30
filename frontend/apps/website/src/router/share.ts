import { RouteRecordRaw } from 'vue-router'
import MainLayout from '@/layouts/MainLayout.vue'

/**
 * 项目密钥校验
 */
export const projectVerificationRoute: RouteRecordRaw = {
  name: 'proejct.verification.route',
  path: 'project/:project_id/verification',
  meta: { ignoreAuth: true },
  component: MainLayout,
}

/**
 * 文档密钥校验
 */
export const collectionVerificationRoute: RouteRecordRaw = {
  name: 'collection.verification.route',
  path: 'doc/:doc_public_id/verification',
  meta: { ignoreAuth: true },
  component: MainLayout,
}
