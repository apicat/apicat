import {
  getDocumentDetailPath,
  getDocumentEditPath,
  getSchemaDetailPath,
  getSchemaEditPath,
  getDefinitionResponseDetailPath,
  getDefinitionResponseEditPath,
  getDocumentDetailPathWithIterationId,
  getDocumentEditPathWithIterationId,
  getSchemaDetailPathWithIterationId,
  getSchemaEditPathWithIterationId,
  getDefinitionResponseDetailPathWithIterationId,
  getDefinitionResponseEditPathWithIterationId,
} from '@/router'
import { useIterationStore } from '@/store/iteration'
import { useParams } from './useParams'

export const useGoPage = () => {
  const router = useRouter()
  const { project_id, iteration_id, doc_id, schema_id, response_id } = useParams()
  const iterationStore = useIterationStore()

  const goSchemaDetailPage = (schemaId?: string | number, replace?: boolean) => {
    const params = {
      path: getSchemaDetailPath(project_id, schemaId || schema_id),
      replace: replace === undefined ? false : replace,
    }

    if (iterationStore.isIterationRoute) {
      params.path = getSchemaDetailPathWithIterationId(iteration_id, schemaId || schema_id)
    }

    router.push(params)
  }

  const goSchemaEditPage = (schemaId?: string | number, replace?: boolean) => {
    const params = {
      path: getSchemaEditPath(project_id, schemaId || schema_id),
      replace: replace === undefined ? false : replace,
    }

    if (iterationStore.isIterationRoute) {
      params.path = getSchemaEditPathWithIterationId(iteration_id, schemaId || schema_id)
    }

    router.push(params)
  }

  const goResponseDetailPage = (responseId?: string | number, replace?: boolean) => {
    const params = {
      path: getDefinitionResponseDetailPath(project_id, responseId || response_id),
      replace: replace === undefined ? false : replace,
    }

    if (iterationStore.isIterationRoute) {
      params.path = getDefinitionResponseDetailPathWithIterationId(iteration_id, responseId || response_id)
    }

    router.push(params)
  }

  const goResponseEditPage = (responseId?: string | number, replace?: boolean) => {
    const params = {
      path: getDefinitionResponseEditPath(project_id, responseId || response_id),
      replace: replace === undefined ? false : replace,
    }

    if (iterationStore.isIterationRoute) {
      params.path = getDefinitionResponseEditPathWithIterationId(iteration_id, responseId || response_id)
    }

    router.push(params)
  }

  const goDocumentDetailPage = (docId?: string | number, replace?: boolean) => {
    const params = {
      path: getDocumentDetailPath(project_id, docId || doc_id),
      replace: replace === undefined ? false : replace,
    }

    if (iterationStore.isIterationRoute) {
      params.path = getDocumentDetailPathWithIterationId(iteration_id, docId || doc_id)
    }

    router.push(params)
  }

  const goDocumentEditPage = (docId?: string | number, replace?: boolean) => {
    const params = {
      path: getDocumentEditPath(project_id, docId || doc_id),
      replace: replace === undefined ? false : replace,
    }

    if (iterationStore.isIterationRoute) {
      params.path = getDocumentEditPathWithIterationId(iteration_id, docId || doc_id)
    }

    router.push(params)
  }

  return {
    goSchemaDetailPage,
    goSchemaEditPage,

    goDocumentDetailPage,
    goDocumentEditPage,

    goResponseDetailPage,
    goResponseEditPage,
  }
}
