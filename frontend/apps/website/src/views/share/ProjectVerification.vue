<template>
  <VerificationForm :handle-check-secret-key="handleCheckSecretKey" />
</template>

<script setup lang="ts">
import VerificationForm from './components/VerificationForm.vue'
import { checkProjectSecret, setProjectSharedToken } from '@/api/project'

const { params } = useRoute()
const handleCheckSecretKey = async (secret_key: string) => {
  try {
    const { project_id } = params
    const { token } = await checkProjectSecret({ project_id, secret_key })
    setProjectSharedToken(project_id as string, token)
  } catch (error) {
    //
  }
}
</script>
