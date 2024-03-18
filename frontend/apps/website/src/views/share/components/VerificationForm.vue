<script setup lang="ts">
import { useNamespace } from '@apicat/hooks'

interface Props {
  handleCheckSecretKey: (secretKey: string) => Promise<any>
}

const props = defineProps<Props>()
const ns = useNamespace('verification')
const isLoading = ref(false)
const form = reactive({
  secretKey: '',
})

async function onSubmitBtnClick() {
  isLoading.value = true
  await props.handleCheckSecretKey(form.secretKey)
  isLoading.value = false
}
</script>

<template>
  <main :class="ns.b()">
    <div :class="ns.e('main')">
      <AcLogo />
      <div class="w-1/2 my-7">
        <el-input v-model="form.secretKey" :placeholder="$t('app.verifyShare.holder')" maxlength="6" clearable />
      </div>
      <el-button :loading="isLoading" type="primary" @click="onSubmitBtnClick">
        {{ $t('app.verifyShare.btn') }}
      </el-button>
      <img src="@/assets/images/img-join.png" class="w-full mt-9" />
    </div>
  </main>
</template>

<style lang="scss" scoped>
@use '@/styles/mixins/mixins' as *;

@include b(verification) {
  @apply h-screen w-screen flex justify-center;

  @include e(main) {
    @apply bg-white shadow-2xl flex flex-col items-center pt-6 h-fit fixed top-56 rounded w480px overflow-hidden;
  }
}
</style>
