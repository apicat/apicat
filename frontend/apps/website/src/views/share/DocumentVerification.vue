<template>
  <VerificationForm :handle-check-secret-key="handleCheckSecretKey" />
</template>

<script setup lang="ts">
import uesShareStore from '@/store/share'
import VerificationForm from './components/VerificationForm.vue'
import { checkCollectionSecret, setCollectionSharedToken } from '@/api/collection'
const shareStore = uesShareStore()
const { params } = useRoute()

const handleCheckSecretKey = async (secret_key: string) => {
  if (!shareStore.sharedDocumentInfo) {
    return
  }

  try {
    const { project_id, collection_id } = shareStore.sharedDocumentInfo
    const token = await checkCollectionSecret({ project_id, collection_id, secret_key })
    params.doc_public_id && setCollectionSharedToken(params.doc_public_id as string, token)
  } catch (error) {
    //
  }
}
</script>
