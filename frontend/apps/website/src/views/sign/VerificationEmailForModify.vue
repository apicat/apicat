<script setup lang="ts">
import { BadRequestError } from '@/api/error'
import { changeUserEmail } from '@/api/sign/user'
import { useInitedPageWithGlobalLoading } from '@/hooks'
import { MAIN_PATH, NOT_FOUND_PATH } from '@/router'

const result = ref<CommonResponseMessageForMessageTemplate>()
const hasError = ref(false)
const router = useRouter()

onBeforeMount(
  useInitedPageWithGlobalLoading(async () => {
    const code = useRoute().params.token as string

    try {
      result.value = (await changeUserEmail(code)) as any
      hasError.value = false
    } catch (error) {
      hasError.value = true
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
  <MessageTemplate v-bind="result">
    <template v-if="!hasError" #description>
      <AcCountDown v-slot="{ seconds }" :time="5" :link="MAIN_PATH" auto-jump>
        <p>
          {{ $t('app.verifyEmail.p1') + seconds + $t('app.verifyEmail.p2')
          }}<RouterLink class="text-primary" :replace="true" :to="MAIN_PATH">
            {{ $t('app.verifyEmail.p3') }} </RouterLink
          >.
        </p>
      </AcCountDown>
    </template>
  </MessageTemplate>
</template>
