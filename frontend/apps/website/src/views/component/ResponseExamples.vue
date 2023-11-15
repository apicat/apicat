<template>
  <div class="relative" v-if="!hasExamples">
    <div class="flex justify-between mt-10px">
      <p> {{ $t('app.response.tips.responseExample') }}<el-tag v-if="lang" class="ml-4px" disable-transitions effect="plain">format:{{ lang }}</el-tag></p>
      <el-button v-if="!readonly" @click="handleAddExample" link type="primary"><el-icon><ac-icon-ep-plus /></el-icon>添加示例</el-button>
    </div>
    <el-tabs v-model="activeTabName" :closable="!readonly" @tab-remove="handleRemoveExample">
      <el-tab-pane :label="item.summary || 'Example ' + (index + 1)" v-for="(item, index) in examples" :key="index" :name="index">
        <p v-if="!readonly" class="mb-10px">示例名称</p>
        <el-input v-if="!readonly" :model-value="item.summary" :readonly="readonly" @update:model-value="(v) => changeExampleFeild(item, 'summary', v)" maxlength="150" placeholder="示例名称" />
        <p v-if="!readonly" class="my-10px">示例值</p>
        <CodeEditor :model-value="item.value" :readonly="readonly" @update:model-value="(v) => changeExampleFeild(item, 'value', v)" :lang="lang" />
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script lang="ts" setup>
import CodeEditor from '@/components/APIEditor/CodeEditor.vue'
import { ref } from 'vue'
import debounce from 'lodash-es/debounce'
interface Example {
  summary: string;
  value: string
}

type Examples = Record<string, Example>

interface Props {
  examples: Examples
  readonly?:boolean
  lang?: string
}

interface Emits {
  (e: 'update:examples', examples: Examples): void
}

const props = defineProps<Props>()
const emits = defineEmits<Emits>()
const examples = ref<Example[]>(initExamples())
const activeTabName = ref(0)

function notify() {
  emits('update:examples', examples.value.reduce((acc, cur, index) => {
    const key = index.toString()
    acc[key] = cur
    return acc
  }, {} as any))
}

const debounceNotify = debounce(notify, 200)

function initExamples(obj: Examples = props.examples): Example[] {
  const arr: Example[] = []
  Object.keys(obj).forEach(key => {
    const example = { ...(obj[key] || {}) }
    arr.push(example)
  })

  return arr
}

function handleAddExample() {
  const e: Example = {
    summary: '',
    value: ''
  }
  examples.value.push(e)
  activeTabName.value = examples.value.length - 1
  notify()
}

function handleRemoveExample(index: any) {
  examples.value.splice(index, 1)
  if(index === activeTabName.value){
    activeTabName.value = examples.value.length - 1
  }

  notify()
}

function changeExampleFeild(e: any, feild: string, v: string) {
  e[feild] = v
  debounceNotify()
}

const hasExamples = computed(()=> props.readonly && !examples.value.length)

watch(() => props.examples, () => {
  examples.value = initExamples()
})
</script>
