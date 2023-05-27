<template>
  <div :class="[ns.b(), ns.is('readonly')]">
    <label v-if="nodeAttrs.method" :class="ns.e('method')" :style="{ backgroundColor: getRequestMethodColor(nodeAttrs.method) }">
      {{ nodeAttrs.method.toUpperCase() }}
    </label>

    <div :class="ns.e('server')">
      <span class="copy_text">{{ mockServerPath }}</span>
    </div>
    <p :class="ns.e('path')" class="copy_text">{{ pathRef }}</p>
    <el-tooltip :content="$t('app.common.fetchMockData')">
      <i :class="ns.e('copy')" @click="handlerMock(fullPath)">
        <ac-icon-quill:send />
      </i>
    </el-tooltip>
    <el-tooltip :content="$t('app.common.copyAllPath')">
      <i :class="ns.e('copy')" class="copy_text" :data-text="fullPath">
        <ac-icon-ic-outline-content-copy />
      </i>
    </el-tooltip>
  </div>
</template>
<script setup lang="ts">
import { useNamespace } from '@/hooks'
import { HttpDocument } from '@/typings'
import { HTTP_URL_NODE_KEY, useNodeAttrs } from '@/hooks/useNodeAttrs'
import { getRequestMethodColor } from '@/commons'
import { mockServerPath, mockApiPath, getMockData } from '@/api/mock'
import { useParams } from '@/hooks/useParams'

const props = defineProps<{ doc: HttpDocument; code: string | number }>()
const { project_id } = useParams()
const ns = useNamespace('http-method')
const nodeAttrs = useNodeAttrs(props, HTTP_URL_NODE_KEY, 'doc')
const pathRef = computed(() => mockApiPath(project_id as string, nodeAttrs.value.path))
const fullPath = computed(() => mockServerPath + pathRef.value)

const handlerMock = async (path: string) => {
  const data = await getMockData(path, { mock_response_code: props.code as string })
  console.log('mock path:', path, ' code:', props.code)
  console.log('mock data:', data)
}
</script>
