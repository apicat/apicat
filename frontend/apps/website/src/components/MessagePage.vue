<script setup lang="ts">
import AheadNav from '@/views/component/AheadNav.vue'

interface OprationMSG {
  emoji?: string
  title?: string
  description?: string
  link?: {
    label: string
    href: string
  }
  countDown?: number
}

const props = defineProps<OprationMSG>()

const defaults = {
  emoji: '',
  title: 'Unknow msg',
  description: '<p>No message provided</p>',
  count: 5,
}

const jump = ref()
const data = computed(() => Object.assign(defaults, props, ((window as any).OPRA_MSG || {})))
const linkInfo = computed(() => {
  if (data.value.link) {
    return {
      label: data.value.link.label,
      href: data.value.link.href,
    }
  }

  return null
})

onMounted(() => {
  if (data.value && data.value.link) {
    if (Object.keys(data.value.link).length > 0) {
      data.value.description.replace('{{link}}', `<a href="${data.value.link.href}">${data.value.link.label}</a>`)
      jump.value = 5
      function wait() {
        if (!jump.value) {
          window.location.href = data.value.link.href
        }
        else {
          jump.value--
          setTimeout(wait, 1000)
        }
      }
      setTimeout(wait, 1000)
    }
  }
})
</script>

<template>
  <AheadNav>
    <div class="flex items-center justify-center w-full text-center">
      <div class="mt-20">
        <span class="text-50px">
          {{ data.emoji }}
        </span>
        <h1 class="mt-5 font-500 text-24px text-gray-title">
          {{ data.title }}
        </h1>
        <div class="mt-5" v-html="data.description" />
      </div>
      <p v-if="linkInfo">
        jump to <a :href="linkInfo.href">Location</a> in {{ jump }} seconds
      </p>
    </div>
  </AheadNav>
</template>
