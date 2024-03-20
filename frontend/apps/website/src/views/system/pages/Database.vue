<script setup lang="ts">
import { ElInput } from 'element-plus'
import { apiGetDatabase } from '@/api/system'

const tBase = 'app.system.db'
const data = ref<SystemAPI.Database>({
  host: '',
  username: '',
  // password: '',
  database: '',
})

apiGetDatabase().then((v) => {
  data.value = v
})
</script>

<template>
  <div class="bg-white w-85%">
    <h1>{{ $t(`${tBase}.title`) }}</h1>
    <ElForm class="content" label-position="top">
      <div style="margin-top: 40px">
        <!-- host -->
        <ElFormItem :label="$t(`${tBase}.host`)">
          <ElInput v-model="data.host" maxlength="255" contenteditable="false" class="h-40px" readonly />
        </ElFormItem>

        <!-- username -->
        <ElFormItem :label="$t(`${tBase}.username`)">
          <ElInput v-model="data.username" maxlength="255" contenteditable="false" class="h-40px" readonly />
        </ElFormItem>

        <!-- <ElFormItem :label="$t(`${tBase}.pw`)"> -->
        <!--   <ElInput maxlength="255" contenteditable="false" v-model="data.password" class="h-40px" readonly /> -->
        <!-- </ElFormItem> -->

        <!-- database -->
        <ElFormItem :label="$t(`${tBase}.db`)">
          <ElInput maxlength="255" contenteditable="false" :model-value="data.database" class="h-40px" readonly />
        </ElFormItem>
      </div>
    </ElForm>
  </div>
</template>

<style scoped>
h1 {
  font-size: 30px;
}

:deep(.el-select .el-input) {
  height: 40px;
}

:deep(.el-button) {
  height: 40px;
}

.row {
  margin-top: 1em;
  margin-bottom: 1em;
  display: flex;
  justify-content: space-between;
  width: 100%;
}
.left,
.right {
  display: flex;
  align-items: center;
}
.left {
  justify-content: flex-start;
  /* flex-grow: 1; */
}
.right {
  /* justify-content: flex-end; */
  flex-grow: 1;
}

.content {
  margin-top: 40px;
}

/* el-upload */
:deep(.content .el-upload) {
  width: 200px;
  height: 200px;
  border-radius: 50%;
}

/* el-image */
.content .block {
  padding: 30px 0;
  text-align: center;
  border-right: solid 1px var(--el-border-color);
  display: inline-block;
  width: 49%;
  box-sizing: border-box;
  vertical-align: top;
}
.content .demonstration {
  display: block;
  color: var(--el-text-color-secondary);
  font-size: 14px;
  margin-bottom: 20px;
}
.content .el-image {
  width: 200px;
  height: 200px;
  border-radius: 50%;
}

.content .image-slot {
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
