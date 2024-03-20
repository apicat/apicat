import { compile } from 'path-to-regexp'
import type { RouteLocationRaw } from 'vue-router'
import { useParams } from './useParams'
import {
  ITERATION_COLLECTION_PATH_NAME,
  ITERATION_RESPONSE_PATH_NAME,
  ITERATION_SCHEMA_PATH_NAME,
  PROJECT_COLLECTION_PATH,
  PROJECT_COLLECTION_PATH_NAME,
  PROJECT_DETAIL_PATH,
  PROJECT_RESPONSE_PATH,
  PROJECT_RESPONSE_PATH_NAME,
  PROJECT_SCHEMA_PATH,
  PROJECT_SCHEMA_PATH_NAME,
} from '@/router'

export function getProjectDetailPath(project_id: number | string) {
  return compile(PROJECT_DETAIL_PATH)({ project_id })
}

export function getCollectionPath(project_id: number | string, collectionID: number | string) {
  return compile(PROJECT_COLLECTION_PATH)({ project_id, collectionID })
}

export function getSchemaPath(project_id: number | string, schemaID: number | string) {
  return compile(PROJECT_SCHEMA_PATH)({ project_id, schemaID })
}

export function getResponsePath(project_id: number | string, responseID: number | string) {
  return compile(PROJECT_RESPONSE_PATH)({ project_id, responseID })
}

export function useGoPage() {
  const router = useRouter()
  const { iterationID } = useParams()
  const goProjectDetailPage = (projectId: string | number, replace?: boolean) => {
    router.push({
      path: getProjectDetailPath(projectId),
      replace: replace === undefined ? false : replace,
    })
  }

  const goCollectionPage = (project_id: string, collectionID: string | number, replace = false) => {
    const r: RouteLocationRaw = {
      replace,
    }
    if (iterationID.value) {
      r.name = ITERATION_COLLECTION_PATH_NAME
      r.params = {
        collectionID,
      }
    }
    else {
      r.name = PROJECT_COLLECTION_PATH_NAME
      r.params = {
        project_id,
        collectionID,
      }
    }

    return router.push(r)
  }

  const goSchemaPage = (project_id: string, schemaID: string | number, replace = false) => {
    const r: RouteLocationRaw = {
      replace,
    }
    if (iterationID.value) {
      r.name = ITERATION_SCHEMA_PATH_NAME
      r.params = {
        schemaID,
      }
    }
    else {
      r.name = PROJECT_SCHEMA_PATH_NAME
      r.params = {
        project_id,
        schemaID,
      }
    }
    return router.push(r)
  }

  const goResponsePage = (project_id: string, responseID: string | number, replace = false) => {
    const r: RouteLocationRaw = {
      replace,
    }
    if (iterationID.value) {
      r.name = ITERATION_RESPONSE_PATH_NAME
      r.params = {
        responseID,
      }
    }
    else {
      r.name = PROJECT_RESPONSE_PATH_NAME
      r.params = {
        project_id,
        responseID,
      }
    }
    return router.push(r)
  }

  return {
    goProjectDetailPage,
    goCollectionPage,
    goSchemaPage,
    goResponsePage,
  }
}
