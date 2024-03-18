<script setup lang="ts">
import VerificationForm from './components/VerificationForm.vue'
import useProjectStore from '@/store/project'

const projectStore = useProjectStore()
async function handleCheckSecretKey(secretCode: string) {
  await projectStore.getProjectShareTokenBySecretCode(secretCode)
}
</script>

<template>
  <VerificationForm :handle-check-secret-key="handleCheckSecretKey" />
</template>
