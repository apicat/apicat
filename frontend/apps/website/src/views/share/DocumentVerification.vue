<script setup lang="ts">
import VerificationForm from './components/VerificationForm.vue'
import { apiSendDocShareKey } from '@/api/project/share'
import { useCollectionsStore } from '@/store/collections'

const props = defineProps<{
  publicID: string
  projectID?: string
  collectionID?: number
}>()
const emit = defineEmits(['update:visible'])
const collectionStore = useCollectionsStore()

// 通过密钥获取collection分享token
async function sendCode(code: string) {
  try {
    const res = await apiSendDocShareKey(props.projectID!, props.collectionID!, code)
    collectionStore.setShareToken(props.publicID, res.shareCode)
    emit('update:visible', true)
  }
  catch (e) {
    //
  }
}
</script>

<template>
  <VerificationForm :handle-check-secret-key="sendCode" />
</template>
