import { RouteRecordRaw } from 'vue-router'
import ProjectDetailLayout from '@/layouts/ProjectDetailLayout'
import { PROJECT_DETAIL_PATH } from './constant'
import { DOCUMENT_DETAIL_NAME, documentDetailRoute, documentEditRoute } from './document'
import { schemaDetailRoute, schemaEditRoute } from './definition.schema'
import { definitionResponseDetailRoute, definitionResponseEditRoute } from './definition.response'
import { compile } from 'path-to-regexp'

export const getProjectDetailPath = (project_id: number | string) => compile(PROJECT_DETAIL_PATH)({ project_id })

export const projectDetailRoute: RouteRecordRaw = {
  name: 'project.detail',
  path: PROJECT_DETAIL_PATH,
  component: ProjectDetailLayout,
  children: [documentDetailRoute, documentEditRoute, schemaDetailRoute, schemaEditRoute, definitionResponseDetailRoute, definitionResponseEditRoute],
  redirect: { name: DOCUMENT_DETAIL_NAME },
  meta: {},
}

export * from './document'
export * from './definition.schema'
export * from './definition.response'
