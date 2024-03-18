<script setup lang="ts">
import { BadRequestError } from '@/api/error'
import { activeAccountByEmail } from '@/api/sign/user'
import { useInitedPageWithGlobalLoading } from '@/hooks'
import { MAIN_PATH, NOT_FOUND_PATH } from '@/router'
import { useUserStore } from '@/store/user'

const result = ref<CommonResponseMessageForMessageTemplate>()
const userStore = useUserStore()
const router = useRouter()

onMounted(
  useInitedPageWithGlobalLoading(async () => {
    const code = useRoute().params.token as string
    try {
      const data = (await activeAccountByEmail(code)) as any
      userStore.updateToken(data.accessToken)
      location.href = MAIN_PATH
    } catch (error) {
      if (error instanceof BadRequestError) {
        const { response } = error as BadRequestError<CommonResponseMessageForMessageTemplate>
        if (!response || (!response.emoji && !response.title && response.message)) router.push(NOT_FOUND_PATH)
        result.value = response
      }
    }
  }),
)
</script>

<template>
  <MessageTemplate v-bind="result" />
</template>
