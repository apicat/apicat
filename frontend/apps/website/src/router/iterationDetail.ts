import type { RouteRecordRaw } from 'vue-router'
import {
  ITERATION_COLLECTION_PATH,
  ITERATION_COLLECTION_PATH_NAME,
  ITERATION_DETAIL_PATH,
  ITERATION_DETAIL_PATH_NAME,
  ITERATION_RESPONSE_PATH,
  ITERATION_RESPONSE_PATH_NAME,
  ITERATION_SCHEMA_PATH,
  ITERATION_SCHEMA_PATH_NAME,
  NOT_FOUND_PATH,
} from './constant'
import { useGlobalLoading } from '@/hooks/useGlobalLoading'
import useProjectStore from '@/store/project'
import { NotFoundError } from '@/api/error'
import { useIterationStore } from '@/store/iteration'
import { useTitle } from '@/hooks/useTitle'

const title = useTitle()

export const iterationDetailRoute: RouteRecordRaw = {
  name: ITERATION_DETAIL_PATH_NAME,
  path: ITERATION_DETAIL_PATH,
  beforeEnter: async (to, _, next) => {
    const { showGlobalLoading, hideGlobalLoading } = useGlobalLoading()
    const iterationStore = useIterationStore()
    const iterationID = to.params.iterationID as string
    showGlobalLoading()

    try {
      const iterationInfo = await iterationStore.getIterationInfo(iterationID)
      const projectID = iterationInfo.project!.id!
      const projectStore = useProjectStore()
      await projectStore.getProjectInfoById(projectID)

      title.value = `${iterationInfo.title || ''} - ${projectStore.project!.title || ''}`
      return next()
    }
    catch (error) {
      hideGlobalLoading()
      // 404 - NotFoundError
      if (error instanceof NotFoundError)
        return next(NOT_FOUND_PATH)

      // default - show error page
      return next(NOT_FOUND_PATH)
    }
  },
  component: async () => import('@/layouts/ProjectDetailLayout/ProjectDetailLayout.vue'),
  children: [
    {
      name: ITERATION_COLLECTION_PATH_NAME,
      path: ITERATION_COLLECTION_PATH,
      component: () => import('@/views/collection/CollectionPage.vue'),
      props: true,
      children: [],
    },
    {
      name: ITERATION_SCHEMA_PATH_NAME,
      path: ITERATION_SCHEMA_PATH,
      component: () => import('@/views/schema/SchemaPage.vue'),
      props: true,
    },
    {
      name: ITERATION_RESPONSE_PATH_NAME,
      path: ITERATION_RESPONSE_PATH,
      component: () => import('@/views/response/ResponsePage.vue'),
      props: true,
    },
  ],
}
