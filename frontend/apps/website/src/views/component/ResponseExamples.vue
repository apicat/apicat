<template>
  <div class="relative" v-if="!hasExamples">
    <div class="flex justify-between mt-10px">
      <p> {{ $t('app.response.tips.responseExample') }}</p>
      <el-button v-if="!readonly" @click="handleAddExample" link type="primary"><el-icon><ac-icon-ep-plus /></el-icon>添加示例</el-button>
    </div>
    <el-tabs v-model="activeTabName" :closable="!readonly" @tab-remove="handleRemoveExample">
      <el-tab-pane :label="item.summary || 'Example ' + (index + 1)" v-for="(item, index) in examples" :key="index" :name="item.id || index">
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
import { markDataWithKey } from '@/commons';
interface Example {
  id?:string
  summary: string;
  value: any
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
const examples = ref<Example[]>([])
const activeTabName = ref<string|number>('')
const localExamples = ref<Examples>({})

function notify() {
   const result = examples.value.reduce((acc, cur, index) => {
    const key = index.toString()
    acc[key] = cur
    return acc
  }, {} as any)

  emits('update:examples', result)
}

const debounceNotify = debounce(notify, 200)

function initExamples(obj: Examples = props.examples): Example[] {
  const arr: Example[] = []
  Object.keys(obj).forEach(key => {
    const item = obj[key] || {}
      markDataWithKey(item,'id')
    // 判断item.value 是不是字符串类型
    if (item.value && typeof item.value !== 'string') {
      item.value = JSON.stringify(item.value || '',null,2)
    }

    arr.push(item)
  })

  return arr
}

function handleAddExample() {
  const e: Example = {
    summary: '',
    value: ''
  }
  markDataWithKey(e,'id')
  examples.value.push(e)
  activeTabName.value = e.id!
  notify()
}

function handleRemoveExample(id: any) {
  const index = examples.value.findIndex((item) => item.id === id)
    examples.value.splice(index, 1)
  if(id === activeTabName.value){
    activeTabName.value = examples.value[examples.value.length - 1].id!
  }
  notify()
}

function changeExampleFeild(e: any, feild: string, v: string) {
  e[feild] = v
  debounceNotify()
}

const hasExamples = computed(()=> props.readonly && !examples.value.length)
let isUpdate = false
watch(() => props.examples, () => {
  localExamples.value = props.examples

  if(!isUpdate){
    isUpdate = true
    examples.value = initExamples()

    if(examples.value.length && !activeTabName.value)
      activeTabName.value = examples.value[0].id!

    nextTick(()=>{
      isUpdate = false
    })
  }
},{
  immediate:true
})
</script>
