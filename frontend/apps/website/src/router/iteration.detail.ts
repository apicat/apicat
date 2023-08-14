import { RouteRecordRaw } from 'vue-router'
import ProjectDetailLayout from '@/layouts/ProjectDetailLayout'
import { ITERATION_DETAIL_PATH, ITERATION_DETAIL_PATH_NAME } from './constant'

// import { DOCUMENT_DETAIL_NAME, documentDetailRoute, documentEditRoute } from './document'
// import { schemaDetailRoute, schemaEditRoute } from './definition.schema'
// import { definitionResponseDetailRoute, definitionResponseEditRoute } from './definition.response'

import { compile } from 'path-to-regexp'

export const getIterationDetailPath = (iteration_id: number | string) => compile(ITERATION_DETAIL_PATH)({ iteration_id })

export const iterationDetailRoute: RouteRecordRaw = {
  name: ITERATION_DETAIL_PATH_NAME,
  path: ITERATION_DETAIL_PATH,
  component: ProjectDetailLayout,
  // children: [documentDetailRoute, documentEditRoute, schemaDetailRoute, schemaEditRoute, definitionResponseDetailRoute, definitionResponseEditRoute],
  // redirect: { name: DOCUMENT_DETAIL_NAME },
}

export * from './document'
export * from './definition.schema'
export * from './definition.response'
