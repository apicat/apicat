<template>
  <el-backtop v-bind="$attrs" class="ac-backtop">
    <el-icon><ac-icon-ep-arrow-up-bold /></el-icon>
  </el-backtop>
</template>

<style lang="scss" scoped>
@use '@/styles/mixins/mixins' as *;

@include b(backtop) {
  background-color: rgba(0, 0, 0, 0.6);
  border-radius: 2px;
  -webkit-box-shadow: 0 1px 3px rgb(0 0 0 / 20%);
  box-shadow: 0 1px 3px rgb(0 0 0 / 20%);
  -webkit-transition: all 0.2s ease-in-out;
  -o-transition: all 0.2s ease-in-out;
  transition: all 0.2s ease-in-out;

  i {
    color: #fff;
    font-size: 20px;
  }
}
</style>
