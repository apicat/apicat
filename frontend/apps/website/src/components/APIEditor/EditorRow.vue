<template>
  <div
    :class="[
      ns.b(),
      {
        [ns.is('ref')]: data.refObj?.name,
        [ns.is('expanded')]: expand,
        [ns.is('const')]: data.label.slice(0, 1) == '<',
      },
    ]"
  >
    <div :class="ns.e('drag')" draggable="true" @dragstart="dragStartHandler" @dragend="dragEndHandler" v-if="!readonly">
      <el-icon><ac-icon-material-symbols-drag-indicator /></el-icon>
    </div>

    <div :class="ns.e('content')" @dragover.prevent="dragOverHandler" @dragleave.prevent="dragLeaveHandler" @drop.prevent="dropHandler">
      <div :class="[ns.e('item'), ns.e('name')]">
        <span :style="{ width: level * 16 + 'px' }" class="indent-spance"></span>
        <el-icon :class="ns.e('expand')" v-if="data.children" @click="toggleExpandHandler"><ac-icon-ep-arrow-right-bold /></el-icon>
        <span v-else class="el-icon"></span>

        <template v-if="!readonly">
          <el-tag disable-transitions v-if="data.label.slice(0, 1) == '<'">{{ data.label.slice(1, -1) }}</el-tag>
          <EditorInput v-else ref="labelInputRef" style="margin-left: -4px" placeholder="name" :value="data.label" :disabled="isRefChildren(data)" @change="changeName" />
        </template>

        <template v-else>
          <el-tag disable-transitions v-if="data.label.slice(0, 1) == '<'">{{ data.label.slice(1, -1) }}</el-tag>
          <span v-else class="copy_text" :title="data.label">{{ data.label }}</span>
        </template>
      </div>

      <div :class="[ns.e('item'), ns.e('type')]">
        <template v-if="!readonly">
          <el-dropdown trigger="click" @visible-change="resetShowRef">
            <el-button text :class="ns.e('ref')" :type="data.refObj ? 'primary' : undefined" :disabled="isRefChildren(data)">
              <div v-if="data.refObj" :class="ns.e('ref_btn')">
                {{ data.refObj.name }}
                <el-tooltip v-if="!data.isSelf" :content="$t('editor.table.removeBinding')" placement="top">
                  <el-icon @click.stop.prevent="unlinkRefHandler"><ac-icon-carbon-unlink /></el-icon>
                </el-tooltip>
              </div>
              <span v-else>{{ data.type }}</span>
            </el-button>
            <template #dropdown>
              <SelectTypeDropmenu :data="data" :show-ref="showRefModal" @showRef="showRefModal = true" @change="changeType" />
            </template>
          </el-dropdown>
        </template>

        <template v-else>
          <el-text v-if="data.refObj" class="truncate" type="primary" :title="data.refObj.name">{{ data.refObj.name }}</el-text>
          <span v-else>{{ data.type }}</span>
        </template>
      </div>

      <div :class="[ns.e('item'), ns.e('required')]">
        <template v-if="!readonly">
          <el-tooltip v-if="data.parent?.type === 'object'" content="required" placement="top" :show-after="368">
            <el-checkbox
              v-if="isRefChildren(data)"
              size="small"
              :disabled="isRefChildren(data)"
              :checked="data.parent?.refObj?.schema.required?.includes(data.label)"
              @change="changeRequired"
            />
            <el-checkbox v-else size="small" :disabled="isRefChildren(data)" :checked="data.parent?.schema.required?.includes(data.label)" @change="changeRequired" />
          </el-tooltip>
        </template>

        <template v-else>
          <span v-if="data.parent?.type === 'object'">
            <span v-if="isRefChildren(data)">{{ data.parent?.refObj?.schema.required?.includes(data.label) ? $t('editor.table.yes') : $t('editor.table.no') }}</span>
            <span v-else="isRefChildren(data)">{{ data.parent?.schema.required?.includes(data.label) ? $t('editor.table.yes') : $t('editor.table.no') }}</span>
          </span>
        </template>
      </div>

      <div :class="[ns.e('item'), ns.e('example')]">
        <template v-if="!readonly">
          <EditorInput
            v-if="['number', 'integer', 'boolean', 'string'].includes(data.type)"
            :placeholder="$t('editor.table.paramExample')"
            :value="data.schema.example"
            :disabled="isRefChildren(data)"
            @change="(v) => changeSchemaField('example', v)"
          />
        </template>

        <template v-else>
          <span v-if="['number', 'integer', 'boolean', 'string'].includes(data.type)" class="truncate copy_text">{{ data.schema.example }}</span>
        </template>
      </div>

      <div :class="[ns.e('item'), ns.e('description')]">
        <template v-if="!readonly">
          <EditorInput
            :placeholder="$t('editor.table.paramDesc')"
            :value="data.schema.description"
            :disabled="isRefChildren(data)"
            @change="(v) => changeSchemaField('description', v)"
          />
        </template>

        <template v-else>
          <span class="copy_text" :title="data.schema.description">{{ data.schema.description }}</span>
        </template>
      </div>

      <div :class="[ns.e('item'), ns.e('mock'), { 'cursor-pointer': !readonly, 'cursor-not-allowed': !isAllowMock(data) }]" @click="mockHandler($event, data)">
        <span class="truncate" :title="data.schema['x-apicat-mock']"> {{ data.schema['x-apicat-mock'] }}</span>
      </div>

      <div :class="[ns.e('item'), ns.e('operation')]" v-if="!readonly">
        <el-button-group size="small">
          <el-tooltip :content="$t('editor.table.addNode')" placement="top" :show-after="368" v-if="data.type === 'object' && !isRefChildren(data) && !data.refObj">
            <el-button text circle @click="addChildHandler">
              <el-icon :size="14"><ac-icon-ep:plus /></el-icon>
            </el-button>
          </el-tooltip>

          <el-popconfirm :title="$t('editor.common.tips.delete')" v-if="!isConstNode(data.label) && !isRefChildren(data)" @confirm="delHandler">
            <template #reference>
              <el-button text circle>
                <el-icon :size="14"><ac-icon-ep-delete /></el-icon>
              </el-button>
            </template>
          </el-popconfirm>
        </el-button-group>
      </div>
    </div>
    <el-collapse-transition>
      <div :class="ns.e('children')" v-if="expand">
        <div :style="intentLineStyle" :class="ns.e('line')"></div>
        <EditorRow :level="level + 1" v-for="item in data.children" :key="item.key" :data="item" :readonly="readonly" />
      </div>
    </el-collapse-transition>
  </div>
</template>

<script setup lang="ts">
import { computed, inject, nextTick, onMounted, ref, shallowRef } from 'vue'
import SelectTypeDropmenu from './SelectTypeDropmenu.vue'
import EditorInput from './EditorInput.vue'
import { JSONSchema, allowMockTypes, constNodeType } from './types'
import type { Tree } from './types'
import { CheckboxValueType, ElMessage } from 'element-plus'
import { useNamespace } from '@/hooks'
import { cloneDeep } from 'lodash-es'
import { mockRulesModal } from '@/components/MockRules'
import { guessMockRule } from '@/components/MockRules/utils'

const props = defineProps<{
  level: number
  data: Tree
  readonly?: boolean
}>()

const ns = useNamespace('schema-row')

const expandsKeys = inject('expandKeys') as Set<string>
const expand = computed(() => expandsKeys.has(props.data.key))
const intentLineStyle = computed(() => {
  let left = props.level * 16 + 6
  return { left: left + 'px' }
})

const toggleExpandHandler = () => {
  if (!props.data.children) {
    return
  }
  expand.value ? expandsKeys.delete(props.data.key) : expandsKeys.add(props.data.key)
}

function isRefChildren(v: Tree): boolean {
  if (v.parent) {
    return v.parent.refObj ? true : isRefChildren(v.parent)
  }
  return false
}

function isAllowMock(data: Tree) {
  if (data.refObj) {
    return false
  }

  if (isRefChildren(data)) {
    return false
  }

  if (!allowMockTypes.includes(data.type)) {
    return false
  }

  return true
}

function isConstNode(v: string) {
  return constNodeType.root == v || constNodeType.items == v
}

const changeNotify = inject('change') as (root?: JSONSchema) => void

const changeRequired = (v: CheckboxValueType) => {
  const sc = props.data.parent?.schema
  if (!sc) {
    return
  }
  if (!sc.required) sc.required = []
  if (v) {
    sc.required.push(props.data.label)
  } else {
    sc.required = sc.required?.filter((v) => v !== props.data.label)
  }
  changeNotify()
}

function resetObject(v: Object) {
  for (let k of Object.keys(v)) {
    if (k != 'description') {
      delete (v as any)[k]
    }
  }
}

const changeType = ({ type, isRef }: any) => {
  const sc = props.data.schema
  resetObject(sc)
  if (!isRef) {
    sc.type = type
    sc['x-apicat-mock'] = guessMockRule({ name: props.data.label, mockType: type })

    if (sc.type == 'array') {
      sc.items = {
        type: 'string',
        'x-apicat-mock': 'string',
      }
    } else if (sc.type == 'object') {
      sc.properties = {}
    }

    // todo array | object mock?
    if (sc.type === 'array' || sc.type === 'object') {
      delete sc['x-apicat-mock']
    }
  } else {
    sc.$ref = type
  }

  changeNotify(props.data.label === constNodeType.root ? sc : undefined)
}

const changeName = (v: string) => {
  if (v === '') {
    return
  }
  const psch = props.data.parent?.schema
  if (psch?.properties) {
    // 单个object下属性名是唯一
    if (psch.properties[v]) {
      ElMessage.error(`参数「${v}」重复`)
      return
    }
    psch.properties[v] = psch.properties[props.data.label]
    // 继承原始必填
    psch.required = psch.required?.map((one) => (one === props.data.label ? v : one))
    // 继承排序
    const orders = psch['x-apicat-orders'] || []
    orders[orders?.indexOf(props.data.label)] = v
    psch['x-apicat-orders'] = orders
    delete psch.properties[props.data.label]
    changeNotify()
  }
}

const changeSchemaField = (path: string, v: string) => {
  ;(props.data.schema as any)[path] = v
  changeNotify()
}

const showRefModal = ref(false)
const resetShowRef = (flag: boolean) => {
  if (flag) {
    showRefModal.value = false
  }
}

const unlinkRefHandler = () => {
  const d = props.data
  if (d.refObj?.schema) {
    const s = cloneDeep(d.refObj.schema)
    if (!s.description) {
      s.description = props.data.schema.description
    }
    if (props.data.label === constNodeType.root) {
      Object.defineProperty(s, '_id', {
        value: d.refObj.schema._id,
        enumerable: false,
        configurable: false,
        writable: false,
      })
      changeNotify(s)
      return
    }

    // unlink object
    if (d.parent?.schema.properties) {
      d.parent.schema.properties[d.label] = s
    }

    // unlink array item
    if (d.parent?.schema.items && (d.parent?.schema.items as any).$ref) {
      delete (d.parent?.schema.items as any).$ref
      Object.assign(d.parent?.schema.items, s)
    }

    changeNotify()
  }
}

const addChildHandler = () => {
  if (!expandsKeys.has(props.data.key)) {
    expandsKeys.add(props.data.key)
  }
  if (props.data.schema.properties) {
    if (props.data.schema.properties['']) {
      return
    }
    props.data.schema.properties[''] = {
      type: 'string',
      'x-apicat-mock': 'string',
    }
    props.data.schema['x-apicat-orders']?.push('')
  }
}

const mockHandler = (e: MouseEvent, row: Tree) => {
  if (props.readonly || !isAllowMock(row)) {
    return
  }

  mockRulesModal.show({
    model: {
      name: row.label,
      mockRule: row.schema['x-apicat-mock'],
      mockType: row.schema.type,
    },
    onOk: (rule: string) => {
      row.schema['x-apicat-mock'] = rule
    },
  })
}

// addchildren auto focus
const labelInputRef = shallowRef()
onMounted(() => {
  if (props.readonly) {
    return
  }
  if (props.data.label === '') {
    nextTick(() => {
      labelInputRef.value.focus()
    })
  }
})

const delHandler = () => {
  const p = props.data.parent?.schema
  if (p) {
    const name = props.data.label
    p['x-apicat-orders'] = p['x-apicat-orders']?.filter((v) => v != name)
    delete p.properties?.[name]
    changeNotify()
  }
}

// drag
const dragDataKey = 'application/apicat-schemanode'
const dragStartHandler = (ev: DragEvent) => {
  if (ev.dataTransfer) {
    ev.dataTransfer.dropEffect = 'move'
    const nodeEle = (ev.target as Element).parentElement
    if (nodeEle) {
      ev.dataTransfer.setDragImage(nodeEle, 0, 0)
      ev.dataTransfer.setData(dragDataKey, props.data.key)
      nextTick(() => {
        nodeEle.classList.add('dragging')
      })
    }
  }
}

function dropTest(ev: DragEvent) {
  ev.preventDefault()
  if (isConstNode(props.data.label) || isRefChildren(props.data)) {
    return 0
  }
  if (ev.dataTransfer?.getData(dragDataKey) == props.data.key) {
    return 0
  }
  const dom = ev.currentTarget as HTMLElement
  if (ev.offsetY < dom.clientHeight / 2) {
    return -1
  }
  return expand && (props.data.type == 'array' || props.data.type == 'object') ? 0 : 1
}
const dragOverHandler = (ev: DragEvent) => {
  const flag = dropTest(ev)
  if (flag === 0) {
    return
  }
  if (ev.dataTransfer) {
    ev.dataTransfer.dropEffect = 'move'
    const dom = ev.currentTarget as HTMLElement
    switch (flag) {
      case -1:
        dom.style.borderBottom = ''
        dom.style.borderTop = '1px var(--primary-color) solid'
        break
      case 1:
        dom.style.borderBottom = '1px var(--primary-color) solid'
        dom.style.borderTop = ''
        break
    }
  }
}

const dragLeaveHandler = (ev: DragEvent) => {
  const dom = ev.currentTarget as HTMLElement
  dom.style.borderBottom = ''
  dom.style.borderTop = ''
}

const dragEndHandler = (ev: DragEvent) => {
  const nodeEle = (ev.target as Element).parentElement
  if (nodeEle) {
    nodeEle.classList.remove('dragging')
  }
}

const fireDropEvent = inject('drop') as (flag: number, to: Tree, source: string) => void

const dropHandler = (ev: DragEvent) => {
  dragLeaveHandler(ev)
  const flag = dropTest(ev)
  if (flag === 0) {
    return
  }
  const dragData = ev.dataTransfer?.getData(dragDataKey)
  if (dragData && dragData != '') {
    fireDropEvent(flag, props.data, dragData)
  }
}
</script>
