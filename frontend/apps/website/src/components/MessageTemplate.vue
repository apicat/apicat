<script setup lang="ts">
import { isEmpty, isObject } from 'lodash-es'
import AheadNav from '@/views/component/AheadNav.vue'

export interface MessageTemplate {
  hideHeader?: boolean
  emoji?: string
  title?: string
  description?: string
  link?: {
    [key: string]: string
  }
}

const props = defineProps<MessageTemplate>()
const description = computed(() => {
  if (!isEmpty(props.link) && isObject(props.link)) {
    let desc = props.description || ''
    Object.keys(props.link).forEach((key) => {
      const text = key
      const href = props.link![key]
      desc = desc.replace(`{${text}}`, `<a class="text-primary" href="${href}">${text}</a>`)
    })

    return desc
  }

  return props.description
})
</script>

<template>
  <div>
    <AheadNav v-if="!hideHeader" />
    <div class="text-center">
      <div class="mt-20">
        <span class="text-50px">
          {{ props.emoji || '' }}
        </span>
        <h1 class="mt-5 font-500 text-24px text-gray-title">
          {{ props.title || '' }}
        </h1>
        <div v-if="!$slots.description" class="mt-5" v-html="description || ''" />
        <div v-else class="mt-5">
          <slot name="description" />
        </div>
      </div>
    </div>
  </div>
</template>
