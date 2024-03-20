<script setup lang="ts">
import { Icon } from '@iconify/vue'

const props = defineProps<{
  avatar: string
  visible: boolean
}>()

const _visible = ref<boolean>(false)
const visible = computed(() => {
  if (props.visible)
    return _visible
  else return false
})
</script>

<template>
  <ElPopover
    :popper-style="{ width: '350px' }"
    :show-arrow="false"
    :visible="visible"
    placement="left"
    trigger="hover"
  >
    <template #reference>
      <slot />
    </template>
    <slot name="content">
      <div v-if="props.avatar" class="w-full">
        <img class="avatar" :src="props.avatar">
        <!-- <img :src="user.avatar" /> -->
      </div>
      <Icon v-else class="avatar-holder" icon="lucide:image" />
    </slot>
  </ElPopover>
</template>

<style scoped>
.avatar-holder {
  border-radius: 50%;
  width: 30px;
  height: 30px;
  background-color: lightgrey;
  color: grey;
  border: 5px lightgrey solid;
}
</style>
