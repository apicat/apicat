<template>
  <div :class="ns.b()">
    <img :class="logoClass" :style="logoStyle" :src="logo" alt="ApiCat" @click="handleClick" />
  </div>
</template>
<script setup>
import { useNamespace } from '@/hooks'
import logoSquare from '@/assets/images/logo-square.svg'
import logoText from '@/assets/images/logo-text.svg'

const porps = defineProps({
  pure: {
    type: Boolean,
    default: false,
  },

  size: {
    type: Number,
    default: 140,
  },

  href: {
    type: String,
    default: '',
  },
})
const ns = useNamespace('logo')
const logo = computed(() => (porps.pure ? logoSquare : logoText))
const logoStyle = computed(() => (porps.pure ? {} : { width: `${porps.size}px` }))
const logoClass = computed(() => [
  ns.e('img'),
  {
    'cursor-pointer': porps.href !== '',
  },
])

const handleClick = () => {
  if (!porps.href) {
    return
  }

  location.href = porps.href
}
</script>
<style lang="scss" scoped>
@use '@/styles/mixins/mixins.scss' as *;

@include b(logo) {
  @apply inline-flex items-center;

  @include e(img) {
    width: 40px;
  }
}
</style>
