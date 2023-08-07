<template>
  <el-dialog v-model="dialogVisible" append-to-body :close-on-click-modal="false" title="Code 导入" destroy-on-close align-center width="600px">
    <el-select v-model="importCodeType" class="w-full mb-10px">
      <el-option v-for="item in codeImportTypes" :key="item.value" :label="item.label" :value="item.value" />
    </el-select>
    <CodeEditor class="w-full h-400px" v-model="code" lang="json" />
    <div slot="footer" class="text-right mt-20px">
      <el-button @click="dialogVisible = false">
        {{ $t('app.common.cancel') }}
      </el-button>
      <el-button type="primary" @click="handleConfirm">
        {{ $t('app.common.confirm') }}
      </el-button>
    </div>
  </el-dialog>
</template>
<script setup lang="ts">
import { useModal } from '@/hooks'
import CodeEditor from '../APIEditor/CodeEditor.vue'
import { convert } from '@/commons/jsonToSchema'
import { ElMessage } from 'element-plus'
import { trim } from 'lodash-es'
const emits = defineEmits(['ok'])

const { dialogVisible, showModel, hideModel } = useModal()
const codeImportTypes = [
  { label: 'JSON', value: 'json' },
  { label: 'JSON Schema', value: 'jsonschema' },
]

const code = ref('')
const importCodeType: Ref<string> = ref(codeImportTypes[0].value)

const handleConfirm = () => {
  try {
    if (!code.value || !trim(code.value)) {
      hideModel()
      return
    }

    const json = JSON.parse(code.value)
    emits('ok', importCodeType.value === 'json' ? convert(json) : json)
    hideModel()
  } catch (error) {
    ElMessage.error('JSON 格式错误')
  }
}
defineExpose({
  show: async () => {
    code.value = ''
    showModel()
  },
})
</script>
