<template>
  <div :class="ns.b()">
    <img :class="ns.e('img')" src="@/assets/images/logo-apicat.png" alt="ApiCat" />
    <span v-if="!pure" :class="ns.e('text')">ApiCat</span>
  </div>
</template>
<script setup>
defineProps({
  pure: {
    type: Boolean,
    default: false,
  },
})
import { useNamespace } from '@/hooks'
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
    width: 46px;
  }

  @include e(text) {
    font-size: 24px;
    color: rgb(36, 36, 36);
    letter-spacing: 0;
    max-width: 188px;
    text-overflow: ellipsis;
    white-space: nowrap;
    overflow: hidden;
    font-weight: 400;
    font-family: AVANTGA, sans-serif;
    margin-top: 8px;
    margin-left: 10px;
  }
}
</style>
