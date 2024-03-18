<template>
  <div>
    <Icon class="avatar-holder" icon="lucide:image" />
  </div>
</template>

<script setup lang="ts">
import { Icon } from '@iconify/vue'
</script>

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
