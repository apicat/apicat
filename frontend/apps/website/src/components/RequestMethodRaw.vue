<template>
  <div :class="[ns.b(), ns.is('readonly')]">
    <label v-if="nodeAttrs.method" :class="ns.e('method')" :style="{ backgroundColor: methodBgColor }">
      {{ nodeAttrs.method.toUpperCase() }}
    </label>

    <div :class="ns.e('server')" v-if="urls.length">
      <el-dropdown trigger="click" @command="onSelectUrl" placement="bottom-start">
        <label>
          <el-icon class="mr-4px"><ac-icon-ep-caret-bottom /></el-icon>
        </label>
        <template #dropdown>
          <el-dropdown-menu>
            <el-dropdown-item :key="item.url + index" :command="item.url" v-for="(item, index) in urls">{{ item.url }}</el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>

      <span class="copy_text">{{ currentUrl }}</span>
    </div>
    <p :class="ns.e('path')" class="copy_text">{{ nodeAttrs.path }}</p>
    <i :class="ns.e('copy')" class="copy_text" :data-text="currentUrl + nodeAttrs.path">
      <ac-icon-ic-outline-content-copy />
    </i>
  </div>
</template>
<script setup lang="ts">
import { useNamespace } from '@/hooks'
import { HttpDocument } from '@/typings'
import { HTTP_URL_NODE_KEY, useNodeAttrs } from '@/hooks/useNodeAttrs'
import { HttpMethodTypeMap } from '@/commons'

const props = defineProps<{ doc: HttpDocument; urls: Array<any> }>()

const ns = useNamespace('http-method')
const nodeAttrs = useNodeAttrs(props, HTTP_URL_NODE_KEY, 'doc')
const httpUrl = unref(nodeAttrs)
const methodBgColor = computed(() => (HttpMethodTypeMap as any)[httpUrl.method].color)

const currentUrl = ref('')

watch(
  () => props.urls,
  () => {
    currentUrl.value = props.urls && props.urls.length ? props.urls[0].url : ''
  },
  {
    immediate: true,
    deep: true,
  }
)

const onSelectUrl = (url: string) => {
  currentUrl.value = url
}
</script>
