import { ref } from 'vue'
import { apiGetCollectDetail } from '@/api/project/collection'
import { useCollectionContextWithoutMounted } from '@/hooks/useCollectionContext'
import { apiGetDocShareStatus } from '@/api/project/share'
import { NOT_FOUND_PATH } from '@/router'
import { NotFoundError, UnauthorizedError } from '@/api/error'
import useApi from '@/hooks/useApi'
import { useCollectionsStore } from '@/store/collections'

export function useCollectionShare() {
  const route = useRoute()
  const router = useRouter()
  const collectionStore = useCollectionsStore()

  const publicID = route.params.collectionPublicID as string

  collectionStore.initShareToken(publicID)

  const hideVerification = ref(!!collectionStore.shareToken)
  const collectionStatus = ref<ShareAPI.ResponseDocShareStatus>()
  const collectionInfo = ref<CollectionAPI.ResponseCollectionDetail>()
  const { context, initContextData } = useCollectionContextWithoutMounted()

  async function _getShareStatus() {
    try {
      const res = await apiGetDocShareStatus(publicID)
      collectionStatus.value = res
      return res
    }
    catch (e) {
      if (e instanceof NotFoundError)
        router.replace(NOT_FOUND_PATH)
    }
  }

  async function getDocument() {
    const { projectID, collectionID } = collectionStatus.value!
    collectionInfo.value = await apiGetCollectDetail(projectID!, collectionID!)
  }

  async function _getCollectionData() {
    try {
      return await Promise.all([getDocument(), initContextData(collectionStatus.value!.projectID!)])
    }
    catch (e) {
      if (e instanceof UnauthorizedError)
        hideVerification.value = false
    }
  }

  const [statusLoading, getShareStatus] = useApi(_getShareStatus)
  const [dataLoading, getCollectionData] = useApi(_getCollectionData)

  async function onVerifyCodeInputSuccess() {
    await getCollectionData()
  }

  onBeforeMount(async () => {
    if (!(await getShareStatus()))
      return
    if (hideVerification.value)
      await getCollectionData()
  })

  return {
    ...context,
    collectionInfo,
    collectionStatus,
    publicID,

    hideVerification,
    loading: computed(() => statusLoading.value || dataLoading.value),
    onVerifyCodeInputSuccess,
  }
}
