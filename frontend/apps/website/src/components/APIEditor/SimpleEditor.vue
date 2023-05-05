<template>
  <div class="ac-sce-simple">
    <table class="w-full table-fixed readonly" v-if="readonly">
      <tr>
        <th style="width: 38%">{{ $t('editor.table.paramName') }}</th>
        <th style="width: 150px">{{ $t('editor.table.paramType') }}</th>
        <th class="text-center" style="width: 80px">{{ $t('editor.table.required') }}</th>
        <th style="width: 38%">{{ $t('editor.table.paramExample') }}</th>
        <th style="width: 38%">{{ $t('editor.table.paramDesc') }}</th>
      </tr>

      <slot name="before" />

      <tr v-for="(data, index) in list" :key="index">
        <td>
          <span class="break-all copy_text">{{ data.name }}</span>
        </td>
        <td>
          {{ data.schema.type }}
        </td>
        <td class="text-center">
          {{ data.required ? $t('editor.table.yes') : $t('editor.table.no') }}
        </td>
        <td>
          <span class="copy_text">{{ data.schema.example }}</span>
        </td>
        <td class="break-all">
          {{ data.schema.description }}
        </td>
      </tr>
    </table>
    <table class="w-full table-fixed" v-else>
      <tr @dragover="dragOverHandler($event, -1)" @dragleave="dragLeaveHandler" @drop="dropHandler($event, -1)">
        <th class="text-center" style="width: 1px" v-show="draggable"></th>
        <th style="width: 34%">{{ $t('editor.table.paramName') }}</th>
        <th style="width: 150px">{{ $t('editor.table.paramType') }}</th>
        <th class="text-center" style="width: 80px">{{ $t('editor.table.required') }}</th>
        <th style="width: 34%">{{ $t('editor.table.paramExample') }}</th>
        <th style="width: 38%">{{ $t('editor.table.paramDesc') }}</th>
        <th class="text-center" style="width: 50px"></th>
      </tr>
      <tbody>
        <slot name="before" />

        <tr v-for="(data, index) in list" :key="index" @dragover="dragOverHandler($event, index)" @dragleave="dragLeaveHandler" @drop="dropHandler($event, index)">
          <td class="text-center" @dragstart="dragStartHandler($event, index)" @dragend="dragEndHandler" :draggable="draggable" v-show="draggable">
            <el-icon class="mt-5px">
              <ac-icon-material-symbols-drag-indicator />
            </el-icon>
          </td>
          <td>
            <el-input v-model="data._name" @input="(v) => onParamNameChange(data, v)" />
          </td>
          <td class="text-center">
            <el-select v-model="data.schema.type" @change="changeNotify(data)">
              <el-option v-for="item in ['string', 'integer', 'number', 'array', 'boolean']" :key="item" :label="item" :value="item" />
            </el-select>
          </td>

          <td class="text-center">
            <el-checkbox size="small" v-model="data.required" @change="changeNotify(data)" tabindex="0" />
          </td>

          <td>
            <el-input v-model="data.schema.example" @input="changeNotify(data)" />
          </td>
          <td>
            <el-input v-model="data.schema.description" @input="changeNotify(data)" />
          </td>
          <td class="text-center">
            <slot name="operate" :row="data" :index="index" :delHandler="delHandler">
              <el-popconfirm title="delete this?" @confirm="delHandler(index)">
                <template #reference>
                  <el-button size="small" text circle tabindex="-1">
                    <el-icon :size="14">
                      <ac-icon-ep-delete />
                    </el-icon>
                  </el-button>
                </template>
              </el-popconfirm>
            </slot>
          </td>
        </tr>

        <tr>
          <td v-show="draggable"></td>
          <td>
            <el-input v-model="newname" :placeholder="$t('editor.table.addParam')" @keyup.enter="addHandler(newname)">
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
import type { APICatSchemaObject } from './types'
import { useSchemaList } from './useSchemaList'

const props = withDefaults(
  defineProps<{
    readonly?: boolean
    draggable?: boolean
    modelValue: APICatSchemaObject[]
    onChange?: (v: APICatSchemaObject) => void
    onCreate?: (v: APICatSchemaObject) => void
    onDelete?: (v: APICatSchemaObject) => void
  }>(),
  {
    readonly: false,
    draggable: true,
    modelValue: () => [],
  }
)

const emits = defineEmits(['update:modelValue'])

const { newname, model, delHandler, addHandler, onParamNameChange, changeNotify } = useSchemaList(props, emits, (models) =>
  models.map((item) => {
    const newItem = toRaw({ ...item })
    delete newItem._name
    return newItem
  })
)

const list = computed(() => {
  model.value = (props.modelValue || []).map((item: APICatSchemaObject) => ({ ...item, _name: item.name }))
  return model.value
})

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
  if (props.readonly || !props.draggable) {
    return
  }
  ev.preventDefault()
  if (ev.dataTransfer?.getData(dragKey) == i.toString()) return
  const dom = ev.currentTarget as HTMLElement
  dom.classList[ev.offsetY > dom.clientHeight / 2 ? 'add' : 'remove']('drop-indicator')
}

const dragLeaveHandler = (ev: DragEvent) => {
  const dom = ev.currentTarget as HTMLElement
  dom.classList.remove('drop-indicator')
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
    model.value.splice(i + 1, 0, model.value.splice(p, 1)[0])
    changeNotify()
  }
}
</script>
