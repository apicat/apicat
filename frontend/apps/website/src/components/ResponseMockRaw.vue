<template>
  <div :class="[ns.b(), ns.is('readonly')]">
    <label :class="ns.e('method')" :style="{ backgroundColor: getRequestMethodColor(nodeAttrs.method) }"> Mock </label>

    <div :class="ns.e('server')">
      <span class="copy_text">{{ mockServerPathRef }}</span>
    </div>
    <p :class="ns.e('path')" class="copy_text">{{ nodeAttrs.path }}</p>

    <el-tooltip :content="$t('app.common.fetchMockData')">
      <i :class="ns.e('copy')" v-if="!isFetchMockData" @click="handlerMock(fullPath, nodeAttrs.method)">
        <ac-icon-quill:send />
      </i>
      <i :class="ns.e('copy')" v-if="isFetchMockData">
        <ac-icon-ep-loading class="animate-spin" />
      </i>
    </el-tooltip>
    <el-tooltip :content="$t('app.common.copyAllPath')">
      <i :class="ns.e('copy')" class="copy_text" :data-text="fullPath">
        <ac-icon-ic-outline-content-copy />
      </i>
    </el-tooltip>
  </div>
</template>
<script setup lang="tsx">
import { useNamespace } from '@/hooks'
import { HttpDocument } from '@/typings'
import { HTTP_URL_NODE_KEY, useNodeAttrs } from '@/hooks/useNodeAttrs'
import { getRequestMethodColor } from '@/commons'
import { mockServerPath, mockApiPath, getMockData } from '@/api/mock'
import { useParams } from '@/hooks/useParams'
import { AsyncMsgBox } from './AsyncMessageBox'
import { CodeEditor } from './APIEditor'

const props = defineProps<{ doc: HttpDocument; code: string | number }>()
const { project_id } = useParams()
const ns = useNamespace('http-method')
const nodeAttrs = useNodeAttrs(props, HTTP_URL_NODE_KEY, 'doc')
const mockServerPathRef = computed(() => mockServerPath + mockApiPath(project_id as string))
const fullPath = computed(() => mockServerPathRef.value + nodeAttrs.value.path)

const isFetchMockData = ref(false)

const handlerMock = async (path: string, method: string) => {
  isFetchMockData.value = true
  try {
    const data: any = await getMockData(path, method, { mock_response_code: props.code as string })

    AsyncMsgBox({
      title: 'Mock Data',
      width: '50vw',
      showCancelButton: false,
      showConfirmButton: false,
      customStyle: { '--el-messagebox-width': '50vw' },
      message: () => <CodeEditor modelValue={JSON.stringify(data, null, 2)} lang="json" readonly />,
    })
  } catch (error) {
    //
  } finally {
    isFetchMockData.value = false
  }
}
</script>
