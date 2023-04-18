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
          <el-text>{{ modelValue.required?.includes(data.name) ? '是' : '否' }}</el-text>
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
            <el-input v-model="data._name" @input="(v) => onParamNameChange(data, v)" />
          </td>
          <td class="text-center">
            <el-select v-model="data.schema.type" @change="changeNotify">
              <el-option v-for="item in schemaType" :key="item" :label="item" :value="item" />
            </el-select>
          </td>

          <td class="text-center">
            <el-tooltip content="required" placement="top" :show-after="368">
              <el-checkbox size="small" :checked="modelValue.required?.includes(data.name)" @change="(v) => onRequiredChange(data, v)" />
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
import { useSchemaList } from './useSchemaList'

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

const emits = defineEmits(['update:modelValue'])

const transformModel = (models: any) => {
  const ps: Record<string, JSONSchema> = {}

  models.forEach((ele: any) => {
    ps[ele.name] = toRaw(ele.schema)
  })

  return {
    type: 'object',
    properties: ps,
    required: props.modelValue.required,
    'x-apicat-orders': models.map((a: any) => a.name),
  }
}

const onParamNameValid = (oldName: string, newName: string) => {
  const schema = props.modelValue as JSONSchema
  schema.required = schema.required?.map((one) => (one === oldName ? newName : one))
}

const { newname, model, delHandler, addHandler, onParamNameChange, changeNotify } = useSchemaList(emits, transformModel, onParamNameValid)

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
      arr.push({ name: k, _name: k, schema: ps[k] })
    }
    schema['x-apicat-orders'] = orders
  }

  model.value = arr

  return model.value
})

const onRequiredChange = (item: any, isChecked: any) => {
  const schema = props.modelValue as JSONSchema

  if (!schema.required) schema.required = []

  if (isChecked) {
    schema.required.push(item.name)
  } else {
    schema.required = schema.required.filter((v) => v !== item.name)
  }
  schema.required = Array.from(new Set(schema.required))
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
