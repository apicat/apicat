<template>
  <header :class="ns.b()">
    <AcLogo :size="80" href="/" />
  </header>
</template>

<script setup lang="ts">
import { useNamespace } from '@/hooks'
const ns = useNamespace('header')
</script>

<style lang="scss" scoped>
@use '@/styles/mixins/mixins' as *;

@include b(header) {
  @apply fixed top-0 left-0 w-full h48px flex-y-center z-100 px-30px bg-white;
  box-shadow: 0 1px 4px #0015291f;
}
</style>
