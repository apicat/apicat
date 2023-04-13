<template>
  <div :class="ns.b()">
    <el-dropdown trigger="click" @command="handleChooseMethod">
      <label :class="ns.e('method')" :style="{ backgroundColor: methodBgColor }">
        {{ nodeAttrs.method.toUpperCase() }}
        <el-icon class="ml-4px"><ac-icon-ep-arrow-down-bold /></el-icon>
      </label>
      <template #dropdown>
        <el-dropdown-menu class="w-100px">
          <el-dropdown-item :key="key" :command="item" v-for="(item, key) in HttpMethodTypeMap">{{ key.toUpperCase() }}</el-dropdown-item>
        </el-dropdown-menu>
      </template>
    </el-dropdown>
    <input type="text" v-model="nodeAttrs.path" :class="ns.e('path')" :placeholder="placeholder" />
  </div>
</template>
<script setup lang="ts">
import { useNamespace } from '@/hooks'
import { HttpDocument } from '@/typings'
import { HTTP_URL_NODE_KEY, useNodeAttrs } from '@/hooks/useNodeAttrs'
import { HttpMethodTypeMap } from '@/commons'

const props = defineProps<{ modelValue: HttpDocument }>()
const ns = useNamespace('http-method')

const nodeAttrs = useNodeAttrs(props, HTTP_URL_NODE_KEY)

const methodBgColor = computed(() => (HttpMethodTypeMap as any)[nodeAttrs.value.method].color)

const placeholder = 'Path, 以"/"开始'

const handleChooseMethod = (menu: any) => {
  nodeAttrs.value.method = menu.value
}
</script>
