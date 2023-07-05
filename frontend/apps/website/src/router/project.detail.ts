import { RouteRecordRaw } from 'vue-router'
import ProjectDetailLayout from '@/layouts/ProjectDetailLayout'
import { PROJECT_DETAIL_PATH, PROJECT_DETAIL_PATH_NAME } from './constant'
import { DOCUMENT_DETAIL_NAME, documentDetailRoute, documentEditRoute } from './document'
import { schemaDetailRoute, schemaEditRoute } from './definition.schema'
import { definitionResponseDetailRoute, definitionResponseEditRoute } from './definition.response'
import { compile } from 'path-to-regexp'

export const getProjectDetailPath = (project_id: number | string) => compile(PROJECT_DETAIL_PATH)({ project_id })

export const projectDetailRoute: RouteRecordRaw = {
  name: PROJECT_DETAIL_PATH_NAME,
  path: PROJECT_DETAIL_PATH,
  component: ProjectDetailLayout,
  children: [documentDetailRoute, documentEditRoute, schemaDetailRoute, schemaEditRoute, definitionResponseDetailRoute, definitionResponseEditRoute],
  redirect: { name: DOCUMENT_DETAIL_NAME },
  meta: {
    ignoreAuth: true,
  },
}

export * from './document'
export * from './definition.schema'
export * from './definition.response'
