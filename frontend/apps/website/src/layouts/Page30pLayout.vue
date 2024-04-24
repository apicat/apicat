<template>
  <div class="flex flex-col items-center justify-center mx-auto px-36px">
    <div class="text-start min-w-500px md:w-full lg:max-w-800px">
      <slot />
    </div>
  </div>
</template>

<style scoped>
:deep(h1) {
  @apply font-500 border-b border-solid border-gray-lighter pb-15px text-24px text-gray-title mb-20px;
}

:deep(h2) {
  @apply text-gray-title font-300;
}

:deep(.el-select .el-input) {
  height: 40px;
}

:deep(.el-button.is-link) {
  padding: 0;
  vertical-align: unset;
}

:deep(.row) {
  margin-top: 1em;
  margin-bottom: 1em;
  display: flex;
  justify-content: space-between;
  width: 100%;
}

:deep(.left, .right) {
  display: flex;
  align-items: center;
}

:deep(.left) {
  justify-content: flex-start;
}

:deep(.right) {
  flex-grow: 1;
}

:deep(.content) {
  margin-top: 30px;
}

/* el-upload */
:deep(.content .el-upload) {
  width: 200px;
  height: 200px;
  border-radius: 50%;
}

/* el-image */
:deep(.content .block) {
  padding: 30px 0;
  text-align: center;
  border-right: solid 1px var(--el-border-color);
  display: inline-block;
  width: 49%;
  box-sizing: border-box;
  vertical-align: top;
}

:deep(.content .demonstration) {
  display: block;
  color: var(--el-text-color-secondary);
  font-size: 14px;
  margin-bottom: 20px;
}

:deep(.content .el-image) {
  width: 200px;
  height: 200px;
  border-radius: 50%;
}

:deep(.content .image-slot) {
  display: flex;
  justify-content: center;
  align-items: center;
  width: 100%;
  height: 100%;
  background: var(--el-fill-color-light);
  color: var(--el-text-color-secondary);
  font-size: 30px;
}
</style>
