<script lang="ts" setup>
import { computed, nextTick, ref, watch } from 'vue'
import debounce from 'lodash-es/debounce'
import { ElButton, ElIcon, ElInput, ElTabPane, ElTabs } from 'element-plus'
import { CodeMirror, Iconify } from '@apicat/components'
import { useNamespace } from '@apicat/hooks'
import { type WithRowKey, withKeyForRow } from '@apicat/shared'

export type Example = WithRowKey<{
  summary: string
  value: any
}>

export type Examples = Record<string, Example>

export interface ResponseExamplesProps {
  examples: Examples
  readonly?: boolean
  lang?: string
}

export interface ResponseExamplesEmits {
  (e: 'update:examples', examples: Examples): void
}

const props = defineProps<ResponseExamplesProps>()
const emits = defineEmits<ResponseExamplesEmits>()
const examples = ref<Example[]>([])
const activeTabName = ref<string | number>(0)
const localExamples = ref<Examples>({})

function notify() {
  const newObj = examples.value.reduce((acc, cur, index) => {
    const key = index.toString()
    acc[key] = cur
    return acc
  }, {} as Examples)
  emits('update:examples', Object.assign(localExamples.value, newObj))
}

const debounceNotify = debounce(notify, 200)

function initExamples(obj: Examples = props.examples): Example[] {
  const arr: Example[] = []
  Object.keys(obj).forEach((key) => {
    const item = obj[key] || {}

    // 判断item.value 是不是字符串类型
    if (item.value && typeof item.value !== 'string')
      item.value = JSON.stringify(item.value || '', null, 2)

    arr.push(withKeyForRow(item))
  })

  return arr
}

function handleAddExample() {
  const e: Example = withKeyForRow({
    summary: '',
    value: '',
  })
  withKeyForRow(e)
  examples.value.push(e)
  activeTabName.value = e.$rowKey
  notify()
}

function handleRemoveExample(id: any) {
  const index = examples.value.findIndex(item => item.$rowKey === id)
  examples.value.splice(index, 1)
  if (examples.value.length > 0 && id === activeTabName.value)
    activeTabName.value = examples.value[examples.value.length - 1].$rowKey

  notify()
}

function changeExampleFeild(e: any, feild: string, v: string) {
  e[feild] = v
  debounceNotify()
}

const hasExamples = computed(() => props.readonly && !examples.value.length)

let isUpdate = false
watch(
  () => props.examples,
  () => {
    localExamples.value = props.examples

    if (!isUpdate) {
      isUpdate = true
      examples.value = initExamples()
      if (examples.value.length)
        activeTabName.value = examples.value[0].$rowKey
      nextTick(() => {
        isUpdate = false
      })
    }
  },
  {
    immediate: true,
  },
)

const ns = useNamespace('http-response')
</script>

<template>
  <div v-if="!hasExamples" :class="[ns.b(), ns.is('empty', !examples.length)]">
    <div class="mt-10px">
      <p>Examples</p>
    </div>
    <ElTabs v-model="activeTabName" :editable="!readonly" :closable="!readonly" @tab-remove="handleRemoveExample">
      <template #addIcon>
        <ElButton v-if="!readonly" link type="primary" @click="handleAddExample">
          <ElIcon><Iconify icon="ep:plus" /> </ElIcon>Add
        </ElButton>
      </template>
      <ElTabPane
        v-for="(item, index) in examples"
        :key="index"
        :label="item.summary || `Example ${index + 1}`"
        :name="item.$rowKey"
      >
        <p v-if="!readonly" class="m-0 mb-10px">
          Name
        </p>
        <ElInput
          v-if="!readonly"
          :model-value="item.summary"
          :readonly="readonly"
          maxlength="150"
          placeholder="Example Name"
          @update:model-value="(v: any) => changeExampleFeild(item, 'summary', v)"
        />
        <p v-if="!readonly" class="my-10px">
          Example Value
        </p>
        <CodeMirror
          :model-value="item.value"
          :readonly="readonly"
          :lang="lang"
          @update:model-value="(v: any) => changeExampleFeild(item, 'value', v)"
        />
      </ElTabPane>
    </ElTabs>
  </div>
</template>
