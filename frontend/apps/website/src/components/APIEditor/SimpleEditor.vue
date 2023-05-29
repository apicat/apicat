<template>
  <div :class="[nsEditor.b(), { [nsEditor.is('readlony')]: readonly }]">
    <div :class="[nsRow.b(), nsRow.m('header')]">
      <div :class="nsRow.e('content')" @dragover="dragOverHandler($event, -1)" @dragleave="dragLeaveHandler" @drop="dropHandler($event, -1)">
        <div :class="[nsRow.e('item'), nsRow.e('name')]">{{ $t('editor.table.paramName') }}</div>
        <div :class="[nsRow.e('item'), nsRow.e('type')]">{{ $t('editor.table.paramType') }}</div>
        <div :class="[nsRow.e('item'), nsRow.e('required')]">{{ $t('editor.table.required') }}</div>
        <div :class="[nsRow.e('item'), nsRow.e('example')]">{{ $t('editor.table.paramExample') }}</div>
        <div :class="[nsRow.e('item'), nsRow.e('description')]">{{ $t('editor.table.paramDesc') }}</div>
        <div v-if="allowMock" :class="[nsRow.e('item'), nsRow.e('mock')]">{{ $t('editor.table.paramMock') }}</div>
        <div :class="[nsRow.e('item'), nsRow.e('operation')]" v-if="!readonly"></div>
      </div>
    </div>

    <slot name="before"></slot>

    <template v-if="readonly" v-for="(data, index) in list" :key="index">
      <div :class="[nsRow.b()]">
        <div :class="nsRow.e('content')">
          <div :class="[nsRow.e('item'), nsRow.e('name')]">
            <span class="copy_text" :title="data.name">{{ data.name }}</span>
          </div>
          <div :class="[nsRow.e('item'), nsRow.e('type')]">
            {{ data.schema.type }}
          </div>
          <div :class="[nsRow.e('item'), nsRow.e('required')]">
            {{ data.required ? $t('editor.table.yes') : $t('editor.table.no') }}
          </div>
          <div :class="[nsRow.e('item'), nsRow.e('example')]">
            <span class="copy_text" :title="data.schema.example">{{ data.schema.example }}</span>
          </div>
          <div :class="[nsRow.e('item'), nsRow.e('description')]">
            <span class="copy_text" :title="data.schema.description">{{ data.schema.description }}</span>
          </div>
          <div v-if="allowMock" :class="[nsRow.e('item'), nsRow.e('mock')]">
            <span class="truncate" :title="data.schema['x-apicat-mock']"> {{ data.schema['x-apicat-mock'] }}</span>
          </div>
        </div>
      </div>
    </template>

    <template v-else>
      <template v-for="(data, index) in list" :key="index">
        <div :class="[nsRow.b()]" v-if="!data.$ref">
          <div :class="nsRow.e('content')" @dragover="dragOverHandler($event, index)" @dragleave="dragLeaveHandler" @drop="dropHandler($event, index)">
            <div :class="[nsRow.e('item'), nsRow.e('name')]">
              <el-icon v-if="!isEditPath" :class="nsRow.e('drag')" @dragstart="dragStartHandler($event, index)" @dragend="dragEndHandler" :draggable="draggable">
                <ac-icon-material-symbols-drag-indicator />
              </el-icon>
              <el-input v-if="!isEditPath" v-model="data._name" @input="(v) => onParamNameChange(data, v)" />
              <span class="px-12px" v-else>{{ data.name }}</span>
            </div>
            <div :class="[nsRow.e('item'), nsRow.e('type')]">
              <el-select v-if="!isEditPath" v-model="data.schema.type" @change="changeParamType(data)">
                <el-option v-for="item in ['string', 'integer', 'number', 'array', 'boolean']" :key="item" :label="item" :value="item" />
              </el-select>
              <span class="px-11px" v-else>{{ data.schema.type }}</span>
            </div>
            <div :class="[nsRow.e('item'), nsRow.e('required')]">
              <el-checkbox size="small" v-model="data.required" @change="changeNotify(data)" tabindex="0" />
            </div>
            <div :class="[nsRow.e('item'), nsRow.e('example')]">
              <el-input v-model="data.schema.example" @input="changeNotify(data)" />
            </div>
            <div :class="[nsRow.e('item'), nsRow.e('description')]">
              <el-input v-model="data.schema.description" @input="changeNotify(data)" />
            </div>
            <div
              v-if="allowMock"
              :class="[nsRow.e('item'), nsRow.e('mock'), { 'cursor-pointer': !readonly, 'cursor-not-allowed': !isAllowMock(data) }]"
              @click="mockHandler($event, data)"
            >
              <span class="truncate" :title="data.schema['x-apicat-mock']"> {{ data.schema['x-apicat-mock'] }}</span>
            </div>

            <div :class="[nsRow.e('item'), nsRow.e('operation')]">
              <slot v-if="!isEditPath" name="operate" :row="data" :index="index" :delHandler="delHandler">
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
            </div>
          </div>
        </div>
      </template>
      <div :class="[nsRow.b()]" v-if="!isEditPath">
        <div :class="nsRow.e('content')">
          <el-input :class="[nsRow.e('quickly-input')]" v-model="newname" :placeholder="$t('editor.table.addParam')" @keyup.enter="addHandler(newname)">
            <template #suffix>
              <el-icon>
                <ac-icon-mi-enter />
              </el-icon>
            </template>
          </el-input>
        </div>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { useNamespace } from '@/hooks'
import { APICatSchemaObject } from './types'
import { useSchemaList } from './useSchemaList'
import { mockRulesModal } from '@/components/MockRules'

const nsEditor = useNamespace('simple-editor')
const nsRow = useNamespace('simple-row')

const props = withDefaults(
  defineProps<{
    readonly?: boolean
    isEditPath?: boolean
    allowMock?: boolean
    draggable?: boolean
    modelValue: APICatSchemaObject[]
    emptyText?: string
    onChange?: (v: APICatSchemaObject) => void
    onCreate?: (v: APICatSchemaObject) => void
    onDelete?: (v: APICatSchemaObject) => void
  }>(),
  {
    readonly: false,
    isEditPath: false,
    draggable: true,
    modelValue: () => [],
    emptyText: '',
  }
)

const emits = defineEmits(['update:modelValue'])

const { newname, model, delHandler, addHandler, isAllowMock, changeParamType, onParamNameChange, changeNotify } = useSchemaList(props, emits, (models) =>
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

const mockHandler = (e: MouseEvent, row: APICatSchemaObject) => {
  if (props.readonly || !isAllowMock(row)) {
    return
  }

  mockRulesModal.show({
    model: {
      name: row.name,
      mockRule: row.schema['x-apicat-mock'],
      mockType: row.schema.type,
    },
    onOk: (rule: string) => {
      row.schema['x-apicat-mock'] = rule
    },
  })
}

const dragKey = 'application/apicat-sortable'

const dragStartHandler = (ev: DragEvent, i: number) => {
  if (ev.dataTransfer) {
    ev.dataTransfer.dropEffect = 'move'
    const nodeEle = (ev.target as Element).parentElement!.parentElement
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
  const nodeEle = (ev.target as Element).parentElement!.parentElement
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

<style lang="scss">
@use './simple-editor-row.scss';
</style>
