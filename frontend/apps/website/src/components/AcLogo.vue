<template>
  <div :class="ns.b()">
    <img v-if="pure" :class="ns.e('img')" src="@/assets/images/logo.svg" alt="ApiCat" />
    <img v-else :class="[ns.e('img'), ns.is('text')]" src="@/assets/images/logo-text.svg" alt="ApiCat" />
  </div>
</template>
<script setup>
import { useNamespace } from '@/hooks'

defineProps({
  pure: {
    type: Boolean,
    default: false,
  },
})
const ns = useNamespace('logo')
</script>
<style lang="scss" scoped>
@use '@/styles/mixins/mixins.scss' as *;

@font-face {
  font-family: 'AVANTGA';
  src: url('@/assets/font/unicode.avantgd.ttf') format('truetype'), url('@/assets/font/AVANTGA2.ttf') format('truetype');
}

@include b(logo) {
  @apply inline-flex items-center;

  @include e(img) {
    width: 30px;

    @include when(text) {
      width: 140px;
    }
  }
}
</style>
