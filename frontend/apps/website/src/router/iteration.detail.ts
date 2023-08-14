import { RouteRecordRaw } from 'vue-router'
import ProjectDetailLayout from '@/layouts/ProjectDetailLayout'
import { ITERATION_DETAIL_PATH, ITERATION_DETAIL_PATH_NAME } from './constant'

import { ITERATION_DOCUMENT_DETAIL_NAME, iterationDocumentDetailRoute, iterationDocumentEditRoute } from './iteration/document'
import { iterationSchemaDetailRoute, iterationSchemaEditRoute } from './iteration/definition.schema'
import { iterationDefinitionResponseDetailRoute, iterationDefinitionResponseEditRoute } from './iteration/definition.response'

import { compile } from 'path-to-regexp'

export const getIterationDetailPath = (iteration_id: number | string) => compile(ITERATION_DETAIL_PATH)({ iteration_id })

export const iterationDetailRoute: RouteRecordRaw = {
  name: ITERATION_DETAIL_PATH_NAME,
  path: ITERATION_DETAIL_PATH,
  component: ProjectDetailLayout,
  children: [
    iterationDocumentDetailRoute,
    iterationDocumentEditRoute,
    iterationSchemaDetailRoute,
    iterationSchemaEditRoute,
    iterationDefinitionResponseDetailRoute,
    iterationDefinitionResponseEditRoute,
  ],
  redirect: { name: ITERATION_DOCUMENT_DETAIL_NAME },
}

export * from './iteration/document'
export * from './iteration/definition.schema'
export * from './iteration/definition.response'
