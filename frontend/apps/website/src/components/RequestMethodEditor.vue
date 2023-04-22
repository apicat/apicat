<template>
  <div :class="ns.b()">
    <el-dropdown trigger="click" @command="handleChooseMethod">
      <label :class="ns.e('method')" :style="{ backgroundColor: getRequestMethodColor(nodeAttrs.method) }">
        {{ nodeAttrs.method.toUpperCase() }}
        <el-icon class="ml-4px"><ac-icon-ep-arrow-down-bold /></el-icon>
      </label>
      <template #dropdown>
        <el-dropdown-menu class="w-100px">
          <el-dropdown-item :key="key" :command="item" v-for="(item, key) in HttpMethodTypeMap">{{ key.toUpperCase() }}</el-dropdown-item>
        </el-dropdown-menu>
      </template>
    </el-dropdown>
    <input type="text" :value="nodeAttrs.path" @input="onChangePath" :class="ns.e('path')" :placeholder="$t('editor.node.httpMethod.pathPlaceholder')" />
  </div>
</template>
<script setup lang="ts">
import { useNamespace } from '@/hooks'
import { HttpDocument } from '@/typings'
import { HTTP_URL_NODE_KEY, useNodeAttrs } from '@/hooks/useNodeAttrs'
import { HttpMethodTypeMap, getRequestMethodColor } from '@/commons'
import { debounce } from 'lodash-es'
import isURL from 'validator/lib/isURL'
import { ElMessage } from 'element-plus'
import { useI18n } from 'vue-i18n'

const props = defineProps<{ modelValue: HttpDocument }>()

const { t } = useI18n()
const ns = useNamespace('http-method')
const nodeAttrs = useNodeAttrs(props, HTTP_URL_NODE_KEY)

const handleChooseMethod = (menu: any) => {
  nodeAttrs.value.method = menu.value
}

const onChangePath = debounce((e: any) => {
  if (!isURL(e.target.value, { require_valid_protocol: false, require_host: false })) {
    ElMessage.error(t('editor.node.httpMethod.pathError'))
    return
  }
  nodeAttrs.value.path = e.target.value
}, 300)
</script>
