<template>
  <VerificationForm :handle-check-secret-key="handleCheckSecretKey" />
</template>

<script setup lang="ts">
import { getProjectDetailPath } from '@/router/project.detail'
import VerificationForm from './components/VerificationForm.vue'
import { checkProjectSecret, setProjectSharedToken } from '@/api/shareProject'

const { params } = useRoute()
const router = useRouter()

const handleCheckSecretKey = async (secret_key: string) => {
  try {
    const { project_id } = params
    const { token, expiration } = await checkProjectSecret({ project_id, secret_key })
    setProjectSharedToken(project_id as string, token, { expires: expiration })
    router.replace(getProjectDetailPath(project_id as string))
  } catch (error) {
    //
  }
}
</script>
