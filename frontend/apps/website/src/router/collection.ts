import type { RouteRecordRaw } from 'vue-router'
import {
  COLLECTION_SHARE_PATH,
  COLLECTION_SHARE_PATH_NAME,
} from './constant'

import CollectionSharePage from '@/views/collection/CollectionSharePage.vue'

export const collectionShareRoute: RouteRecordRaw = {
  name: COLLECTION_SHARE_PATH_NAME,
  path: COLLECTION_SHARE_PATH,
  component: CollectionSharePage,
  meta: { ignoreAuth: true },
}
