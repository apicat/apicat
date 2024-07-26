<script setup lang="ts">
import { useCollapse } from '@/components/collapse/useCollapse'
import Local from '@/views/system/pages/Storage/Local.vue'
import Cloudflare from '@/views/system/pages/Storage/Cloudflare.vue'
import QiNiu from '@/views/system/pages/Storage/QiNiu.vue'
import { apiGetStorage } from '@/api/system'
import { SysStorage } from '@/commons'

const tBase = 'app.system.storage'
const collapse = useCollapse({ onlyOne: true })

interface A {
  [SysStorage.Disk]: SystemAPI.StorageDisk
  [SysStorage.CF]: SystemAPI.StorageCF
  [SysStorage.Qiniu]: SystemAPI.StorageQiniu
}
const data = ref<A>({
  [SysStorage.Disk]: {
    path: '',
  },
  [SysStorage.CF]: {
    accountID: '',
    accessKeyID: '',
    accessKeySecret: '',
    bucketName: '',
    bucketUrl: '',
  },
  [SysStorage.Qiniu]: {
    accessKey: '',
    secretKey: '',
    bucketName: '',
    bucketUrl: '',
  },
})
apiGetStorage().then((res) => {
  for (let i = 0; i < res.length; i++) {
    const v = res[i]
    data.value[v.driver as keyof A] = v.config as any
    if (v.use)
      collapse.ctx.open(v.driver)
  }
})
</script>

<template>
  <div class="bg-white w-85%">
    <h1>{{ $t(`${tBase}.title`) }}</h1>

    <div class="flex flex-col mt-40px">
      <Local v-model:config="data[SysStorage.Disk]" class="collapse-box" :name="SysStorage.Disk" :collapse="collapse" />
      <Cloudflare
        v-model:config="data[SysStorage.CF]"
        class="collapse-box mt-30px"
        :name="SysStorage.CF"
        :collapse="collapse"
      />
      <QiNiu
        v-model:config="data[SysStorage.Qiniu]"
        class="collapse-box mt-30px"
        :name="SysStorage.Qiniu"
        :collapse="collapse"
      />
    </div>
  </div>
</template>

<style scoped>
h1 {
  font-size: 30px;
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
</style>
