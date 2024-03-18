<script setup lang="ts">
defineProps<{ user: TeamAPI.UserData }>()
</script>

<template>
  <ElPopover :show-after="300" :popper-style="{ width: '350px' }" :show-arrow="false" placement="left" trigger="hover">
    <template #reference>
      <slot />
    </template>
    <div v-if="user.avatar" class="w-full">
      <img :src="user.avatar">
    </div>
  </ElPopover>
</template>

<style scoped>
img {
  width: 350px;
}
</style>
