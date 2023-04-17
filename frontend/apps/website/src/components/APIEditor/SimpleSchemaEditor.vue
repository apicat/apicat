<template>
  <div class="ac-sce-simple">
    <table class="w-full readonly" v-if="readonly">
      <tr>
        <th width="320">参数名</th>
        <th width="138">类型</th>
        <th width="56">必须</th>
        <th width="258">示例值</th>
        <th>描述</th>
      </tr>
      <tr v-for="(data, index) in flatValues" :key="index">
        <td>
          <el-text tag="b">
            <span class="copy_text">{{ data.name }}</span>
          </el-text>
        </td>
        <td>
          <el-text>{{ data.schema.type }}</el-text>
        </td>
        <td>
          <el-text>{{ data.schema.required ? '是' : '否' }}</el-text>
        </td>
        <td>
          <el-text>
            <span class="copy_text">{{ data.schema.example }}</span>
          </el-text>
        </td>
        <td>
          <el-text>{{ data.schema.description }}</el-text>
        </td>
      </tr>
    </table>

    <table class="w-full" v-else>
      <tr @dragover="dragOverHandler($event, -1)" @dragleave="dragLeaveHandler" @drop="dropHandler($event, -1)">
        <th class="text-center" width="32"></th>
        <th>参数名</th>
        <th class="text-center" width="138">类型</th>
        <th class="text-center" width="56">必须</th>
        <th width="258">示例值</th>
        <th width="264">描述</th>
        <th class="text-center" width="48"></th>
      </tr>
      <tbody>
        <tr v-for="(data, index) in flatValues" :key="index" @dragover="dragOverHandler($event, index)" @dragleave="dragLeaveHandler" @drop="dropHandler($event, index)">
          <td class="text-center" @dragstart="dragStartHandler($event, index)" @dragend="dragEndHandler" draggable="true">
            <el-icon class="mt-5px">
              <ac-icon-material-symbols-drag-indicator />
            </el-icon>
          </td>
          <td>
            <el-input v-model="flatValues[index].name" @change="changeNotify" />
          </td>
          <td class="text-center">
            <el-select v-model="data.schema.type" @change="changeNotify">
              <el-option v-for="item in schemaType" :key="item" :label="item" :value="item" />
            </el-select>
          </td>

          <td class="text-center">
            <el-tooltip content="required" placement="top" :show-after="368">
              <el-checkbox size="small" v-model="data.schema.required" @change="changeNotify" />
            </el-tooltip>
          </td>

          <td>
            <el-input v-model="data.schema.default" @change="changeNotify" />
          </td>
          <td>
            <el-input v-model="data.schema.description" @change="changeNotify" />
          </td>
          <td class="text-center">
            <el-popconfirm title="delete this?" @confirm="delHandler(index)">
              <template #reference>
                <el-button text circle>
                  <el-icon :size="14">
                    <ac-icon-ep-delete />
                  </el-icon>
                </el-button>
              </template>
            </el-popconfirm>
          </td>
        </tr>

        <tr>
          <td></td>
          <td>
            <el-input v-model="newname" placeholder="添加参数" @change="addHandler">
              <template #suffix>
                <el-icon>
                  <ac-icon-mi-enter />
                </el-icon>
              </template>
            </el-input>
          </td>
          <td colspan="5"></td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script setup lang="ts">
import { ElMessage } from 'element-plus'
import type { Definition, JSONSchema } from './types'
import { ref, computed } from 'vue'

const emits = defineEmits(['update:modelValue'])

const props = withDefaults(
  defineProps<{
    readonly?: boolean
    modelValue: JSONSchema
    // 引用模型的集合
    definitions?: Definition[]
    hasFile?: boolean
  }>(),
  {
    readonly: false,
    hasFile: false,
    modelValue: () => ({}),
  }
)

const schemaType = ['string', 'integer', 'number', 'array', 'boolean']

if (props.hasFile) {
  schemaType.push('file')
}

const list: any = ref([])

const flatValues = computed(() => {
  const arr: any = []
  const schema = props.modelValue as JSONSchema
  // todo 校验必须是obj？
  let ps = (schema.properties = schema.properties || {})

  if (schema.$ref) {
    const name = schema.$ref.match(/#\/definitions\/(\w+)/)?.[1]
    const refschema = props.definitions?.find((v) => v.name === name)
    if (refschema && refschema.schema) {
      ps = Object.assign({}, refschema.schema.properties)
    }
  }

  if (ps) {
    const orders = schema['x-apicat-orders'] || Object.keys(ps)
    for (let k of orders) {
      // todo additionalMetadata??
      arr.push({ name: k, schema: ps[k] })
    }
    schema['x-apicat-orders'] = orders
  }

  list.value = arr

  return list.value
})

const changeNotify = () => {
  const ps: Record<string, JSONSchema> = {}

  list.value.forEach((ele: any) => {
    ps[ele.name] = toRaw(ele.schema)
  })

  const model = {
    type: 'object',
    properties: ps,
    'x-apicat-orders': list.value.map((a: any) => a.name),
  }

  emits('update:modelValue', model)
}

const newname = ref('')

const addHandler = (v: string) => {
  if (v == '') {
    return
  }

  if (flatValues.value.find((a: any) => a.name == v)) {
    ElMessage.error(`参数「${v}」重复`)
    return
  }
  newname.value = ''
  flatValues.value.push({
    name: v,
    required: false,
    schema: { type: 'string' },
  })

  changeNotify()
}

const delHandler = (i: number) => {
  flatValues.value.splice(i, 1)
  changeNotify()
}

const dragKey = 'application/apicat-sortable'

const dragStartHandler = (ev: DragEvent, i: number) => {
  if (ev.dataTransfer) {
    ev.dataTransfer.dropEffect = 'move'
    const nodeEle = (ev.target as Element).parentElement
    if (nodeEle) {
      nodeEle.classList.add('dragging')
      ev.dataTransfer.setDragImage(nodeEle, 0, 0)
      ev.dataTransfer.setData(dragKey, i.toString())
    }
  }
}

const dragOverHandler = (ev: DragEvent, i: number) => {
  if (props.readonly) {
    return
  }
  ev.preventDefault()
  if (ev.dataTransfer?.getData(dragKey) == i.toString()) return
  const dom = ev.currentTarget as HTMLElement
  dom.style.borderBottom = ev.offsetY > dom.clientHeight / 2 ? '1px blue solid' : ''
}

const dragLeaveHandler = (ev: DragEvent) => {
  const dom = ev.currentTarget as HTMLElement
  dom.style.borderBottom = ''
}

const dragEndHandler = (ev: DragEvent) => {
  const nodeEle = (ev.target as Element).parentElement
  if (nodeEle) {
    nodeEle.classList.remove('dragging')
    nodeEle.style.borderBottom = ''
  }
}

const dropHandler = (ev: DragEvent, i: number) => {
  dragLeaveHandler(ev)
  if (props.readonly) {
    return
  }
  ev.preventDefault()
  const drag = ev.dataTransfer?.getData(dragKey)
  if (drag && drag !== '') {
    const p = parseInt(drag)
    flatValues.value.splice(i + 1, 0, flatValues.value.splice(p, 1)[0])
    changeNotify()
  }
}
</script>

<style>
.ac-sce-simple {
  font-size: 14px;
  border: 1px var(--el-border-color-lighter) solid;
  border-radius: var(--el-border-radius-base);
  overflow: hidden;
}
.ac-sce-simple:focus-within {
  box-shadow: rgba(0, 0, 0, 0.1) 0 2px 4px;
}

.ac-sce-simple table {
  border: 0;
  border-collapse: collapse;
  border-spacing: 0;
}

.ac-sce-simple table th {
  font-weight: normal;
  text-align: left;
  color: var(--el-text-color-secondary);
  background-color: var(--el-fill-color-light);
  padding: 6px 12px;
  cursor: default;
  user-select: none;
  -webkit-user-select: none;
}

.ac-sce-simple table tr {
  border-bottom: 1px var(--el-border-color-lighter) solid;
}

.ac-sce-simple table tr:last-child {
  border-bottom: 0;
}

.ac-sce-simple table tr:hover {
  background-color: var(--el-fill-color-light);
}

.ac-sce-simple table:not(.readonly) tr > td:nth-child(1) {
  opacity: 0.1;
}

.ac-sce-simple table tr:hover > td:nth-child(1) {
  opacity: 1;
}

.ac-sce-simple .el-input__wrapper {
  --el-input-border-color: transparent;
  --el-input-bg-color: transparent;
}

.ac-sce-simple table tr.dragging {
  opacity: 0.3;
}

.ac-sce-simple table.readonly td {
  padding: 6px 12px;
}
</style>
