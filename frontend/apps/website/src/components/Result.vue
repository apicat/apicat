<script setup lang="ts">
interface ResultProps {
  title?: string
  subTitle?: string
}

withDefaults(defineProps<ResultProps>(), {
  title: '',
  subTitle: '',
})
</script>

<template>
  <div class="result">
    <div class="result-icon">
      <slot name="icon" />
    </div>
    <slot name="title">
      <p class="result-title">
        {{ title }}
      </p>
    </slot>

    <p class="mt-1 result-subtitle">
      <slot name="subTitle">
        {{ subTitle }}
      </slot>
    </p>
  </div>
</template>

<style lang="scss" scoped>
.result {
  padding: 48px 32px;
  text-align: center;

  &-icon {
    display: flex;
    justify-content: center;

    img {
      width: 100%;
    }
  }

  &-title {
    @apply text-zinc-600 mt-6px text-24px font-500;
  }

  &-subtitle {
    max-width: 502px;
    margin: auto;
    @apply text-zinc-500;
  }
}
</style>
