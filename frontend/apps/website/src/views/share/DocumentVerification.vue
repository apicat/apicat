<template>
  <VerificationForm :handle-check-secret-key="handleCheckSecretKey" />
</template>

<script setup lang="ts">
import useShareStore from '@/store/share'
import VerificationForm from './components/VerificationForm.vue'
import { checkCollectionSecret, setCollectionSharedToken } from '@/api/collection'
import { getDocumentShareDetailPath } from '@/router/share'
const shareStore = useShareStore()
const { params } = useRoute()
const router = useRouter()

const handleCheckSecretKey = async (secret_key: string) => {
  if (!shareStore.sharedDocumentInfo) {
    return
  }

  try {
    const { project_id, collection_id } = shareStore.sharedDocumentInfo
    const { token, expiration } = await checkCollectionSecret({ project_id, collection_id, secret_key })
    params.doc_public_id && setCollectionSharedToken(params.doc_public_id as string, token, { expires: expiration })
    router.replace(getDocumentShareDetailPath(params.doc_public_id as string))
  } catch (error) {
    //
  }
}
</script>
