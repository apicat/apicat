<template>
  <el-icon v-bind="$attrs" class="ac-icon">
    <slot></slot>
  </el-icon>
</template>
<style lang="scss" scoped>
.ac-icon {
  @apply rounded cursor-pointer hover:bg-gray-110;
}
</style>
