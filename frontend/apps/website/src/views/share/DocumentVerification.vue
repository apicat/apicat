<template>
  <VerificationForm :handle-check-secret-key="handleCheckSecretKey" />
</template>

<script setup lang="ts">
import VerificationForm from './components/VerificationForm.vue'
import { checkCollectionSecret, setCollectionShareInfo } from '@/api/collection'

const { params } = useRoute()
const handleCheckSecretKey = async (secret_key: string) => {
  try {
    const { doc_public_id } = params
    const shareInfo: any = await checkCollectionSecret({ doc_public_id, secret_key })
    setCollectionShareInfo(doc_public_id as string, shareInfo)
  } catch (error) {
    //
  }
}
</script>
