<script setup lang="ts">
import { useCollapse } from '@/components/collapse/useCollapse'
import Local from '@/views/system/pages/Cache/Local.vue'
import Redis from '@/views/system/pages/Cache/Redis.vue'
import { SysCache } from '@/commons'
import { apiGetCache } from '@/api/system'

const tBase = 'app.system.cache'
const collapse = useCollapse({})

interface A {
  [SysCache.Redis]: SystemAPI.CacheRedis
  [SysCache.Local]: SystemAPI.CacheMemory
}
const data = ref<A>({
  [SysCache.Local]: {},
  [SysCache.Redis]: {
    host: '',
    password: '',
    database: 0,
  },
})

const currentUse = ref<SysCache>()
apiGetCache().then((res) => {
  for (let i = 0; i < res.length; i++) {
    const v = res[i]
    data.value[v.driver as keyof A] = v.config as any
    if (v.use) {
      collapse.ctx.open(v.driver)
      currentUse.value = v.driver
    }
  }
})
</script>

<template>
  <div class="bg-white w-85%">
    <h1>{{ $t(`${tBase}.title`) }}</h1>

    <div class="mt-40px flex flex-col">
      <Local
        v-model:config="data[SysCache.Local]"
        class="collapse-box"
        :name="SysCache.Local"
        :collapse="collapse"
        :current-use="currentUse"
      />
      <Redis
        v-model:config="data[SysCache.Redis]"
        class="collapse-box mt-30px"
        :name="SysCache.Redis"
        :collapse="collapse"
        :current-use="currentUse"
      />
    </div>
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
</style>
