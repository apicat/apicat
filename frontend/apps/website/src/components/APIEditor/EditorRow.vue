<template>
  <div class="ac-sce-node" v-if="readonly">
    <div class="ac-sce-node_content" :style="{ paddingLeft: level * 18 + 'px' }">
      <div class="ac-sce-expand" @click="toggleExpandHandler">
        <el-icon v-if="data.children">
          <svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" viewBox="0 0 24 24">
            <path
              d="M12 5.83l2.46 2.46a.996.996 0 1 0 1.41-1.41L12.7 3.7a.996.996 0 0 0-1.41 0L8.12 6.88a.996.996 0 1 0 1.41 1.41L12 5.83zm0 12.34l-2.46-2.46a.996.996 0 1 0-1.41 1.41l3.17 3.18c.39.39 1.02.39 1.41 0l3.17-3.17a.996.996 0 1 0-1.41-1.41L12 18.17z"
              fill="currentColor"
            ></path>
          </svg>
        </el-icon>
      </div>
      <div class="ac-sce-node_body readonly">
        <div>
          <el-space spacer="·">
            <el-tag disable-transitions v-if="data.label.slice(0, 1) == '<'">{{ data.label.slice(1, -1) }}</el-tag>
            <el-text tag="b" v-else>
              <span class="copy_text">{{ data.label }}</span>
              <template v-if="data.parent?.type === 'object'">
                <template v-if="isRefChildren(data)">
                  <el-text v-if="data.parent?.refObj?.schema.required?.includes(data.label)" type="danger"> (*) </el-text>
                </template>
                <template v-else>
                  <el-text v-if="data.parent?.schema.required?.includes(data.label)" type="danger"> (*) </el-text>
                </template>
              </template>
            </el-text>
            <el-text>
              {{ data.type }}
              <el-text v-if="data.refObj" type="primary">({{ data.refObj.name }})</el-text>
            </el-text>
            <el-text class="w-300px" type="info" truncated :title="data.schema.description" v-if="data.schema.description">{{ data.schema.description }}</el-text>
          </el-space>
        </div>
        <div>
          <el-text v-if="data.schema.example" type="info">
            <small>示例</small>: <span class="copy_text">{{ data.schema.example }}</span>
          </el-text>
          <el-text v-if="data.schema.default !== undefined && data.schema.default !== null && data.schema.default !== ''" type="info">
            <small>默认值</small>: <span class="copy_text">{{ data.schema.default }}</span>
          </el-text>
        </div>
      </div>
    </div>
    <el-collapse-transition>
      <div class="ac-sce-node_children">
        <div :style="{ left: level * 18 + 2 + 'px' }" class="indent-line"></div>
        <EditorRow :level="level + 1" v-for="item in data.children" :key="item.key" :data="item" :readonly="readonly" />
      </div>
    </el-collapse-transition>
  </div>
  <div
    v-else
    class="ac-sce-node"
    :aria-expanded="expand"
    :class="{
      'ac-sce-ref': data.refObj?.name,
      expand: expand,
      'ac-sce-const': data.label.slice(0, 1) == '<',
    }"
  >
    <div class="ac-sce-node_drag" draggable="true" @dragstart="dragStartHandler" @dragend="dragEndHandler" v-if="!readonly">
      <el-icon><ac-icon-material-symbols-drag-indicator /></el-icon>
    </div>
    <div
      class="ac-sce-node_content"
      :style="{ paddingLeft: level * 18 + 'px' }"
      @dragover.prevent="dragOverHandler"
      @dragleave.prevent="dragLeaveHandler"
      @drop.prevent="dropHandler"
    >
      <div class="ac-sce-expand" @click="toggleExpandHandler">
        <el-icon v-if="data.children"><ac-icon-ep-arrow-right-bold /></el-icon>
      </div>
      <div class="ac-sce-node_body">
        <div>
          <el-tag disable-transitions v-if="data.label.slice(0, 1) == '<'">{{ data.label.slice(1, -1) }}</el-tag>
          <EditorInput v-else ref="labelInputRef" style="margin-left: -4px" placeholder="name" :value="data.label" :disabled="isRefChildren(data)" @change="changeName" />
        </div>
        <div>
          <el-dropdown trigger="click" @visible-change="resetShowRef">
            <el-button text :type="data.refObj ? 'primary' : undefined" :disabled="isRefChildren(data)">
              <el-space :size="4" v-if="data.refObj">
                {{ data.refObj.name }}
                <el-tooltip content="解除绑定" placement="top">
                  <el-icon @click.stop.prevent="unlinkRefHandler"><ac-icon-carbon-unlink /></el-icon>
                </el-tooltip>
              </el-space>
              <span v-else>{{ data.type }}</span>
            </el-button>
            <template #dropdown>
              <SelectTypeDropmenu :data="data" :show-ref="showRefModal" @showRef="showRefModal = true" @change="changeNotify" />
            </template>
          </el-dropdown>
        </div>
        <div>
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
        </div>
        <div>
          <EditorInput
            v-if="['number', 'integer', 'boolean', 'string'].includes(data.type)"
            placeholder="示例值"
            :value="data.schema.example"
            :disabled="isRefChildren(data)"
            @change="(v) => changeSchemaField('example', v)"
          />
        </div>
        <div>
          <EditorInput placeholder="描述" :value="data.schema.description" :disabled="isRefChildren(data)" @change="(v) => changeSchemaField('description', v)" />
        </div>
        <div>
          <el-button-group size="small">
            <el-tooltip content="添加子节点" placement="top" :show-after="368" v-if="data.type === 'object' && !isRefChildren(data) && !data.refObj">
              <el-button text circle @click="addChildHandler">
                <el-icon :size="14"><ac-icon-ic-outline-add-circle /></el-icon>
              </el-button>
            </el-tooltip>

            <el-popconfirm title="delete this?" v-if="!isConstNode(data.label) && !isRefChildren(data)" @confirm="delHandler">
              <template #reference>
                <el-button text circle>
                  <el-icon :size="14"><ac-icon-iconoir-delete-circle /></el-icon>
                </el-button>
              </template>
            </el-popconfirm>
          </el-button-group>
        </div>
      </div>
    </div>
    <el-collapse-transition>
      <div class="ac-sce-node_children" v-if="expand">
        <div :style="{ left: level * 18 + 2 + 'px' }" class="indent-line"></div>
        <EditorRow :level="level + 1" v-for="item in data.children" :key="item.key" :data="item" />
      </div>
    </el-collapse-transition>
  </div>
</template>

<script setup lang="ts">
import { computed, inject, nextTick, onMounted, ref, shallowRef, unref } from 'vue'
import SelectTypeDropmenu from './SelectTypeDropmenu.vue'
import EditorInput from './EditorInput.vue'
import { JSONSchema, constNodeType } from './types'
import type { Tree } from './types'
import { CheckboxValueType, ElMessage } from 'element-plus'

const props = defineProps<{
  level: number
  data: Tree
  readonly?: boolean
}>()

const expandsKeys = inject('expandKeys') as Set<string>
const expand = computed(() => expandsKeys.has(props.data.key))
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
    const s = Object.assign({}, unref(d.refObj.schema))
    if (!s.description) {
      s.description = props.data.schema.description
    }
    if (props.data.label === constNodeType.root) {
      changeNotify(s)
      return
    } else if (d.parent?.schema.properties) {
      // schema 整体替换的话 需要从上级替换 否则会失去引用
      d.parent.schema.properties[d.label] = s
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
    }
    props.data.schema['x-apicat-orders']?.push('')
  }
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
        dom.style.borderTop = '1px blue solid'
        break
      case 1:
        dom.style.borderBottom = '1px blue solid'
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

<style>
.ac-sce-node {
  user-select: none;
  -webkit-user-select: none;
  font-size: 14px;
  position: relative;
  background-color: var(--el-bg-color);
  transition: opacity 0.3s ease-out;
}

.ac-sce-expand {
  width: 26px;
  height: 26px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: transform 0.15s;
}

.ac-sce-node[aria-expanded='true'] > .ac-sce-node_content > .ac-sce-expand {
  transform: rotate(90deg);
}
.ac-sce-node_content {
  display: flex;
  align-items: center;
  height: 34px;
  position: relative;
  border-width: 1px 0;
  box-sizing: content-box;
  border-style: solid;
  border-color: transparent;
  transition: border 0.3s ease-out;
}
.ac-sce-node_content:hover {
  background-color: var(--el-fill-color-light);
}
.ac-sce-node.dragging {
  opacity: 0.3;
}
.ac-sce-node_drag {
  position: absolute;
  top: 9px;
  left: 6px;
  height: 16px;
  opacity: 0.1;
  z-index: 3;
  /* transform: translateX(-18px); */
}
.ac-sce-node:hover > .ac-sce-node_drag {
  opacity: 1;
}

.ac-sce-node_children {
  position: relative;
}

.indent-line {
  position: absolute;
  top: 0;
  height: 100%;
  margin-left: 10px;
  width: 0px;
  z-index: 2;
  border-left: 1px var(--el-border-color-lighter) dashed;
}

.ac-sce-node:hover > .ac-sce-node_children > .indent-line {
  border-left: 1px var(--el-border-color) solid;
}

.ac-sce-node_body {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: flex-start;
}

.ac-sce-node_body > div:nth-child(1) {
  flex: 1 0 auto;
}
.ac-sce-node_body.readonly > div:nth-child(2) {
  flex: 0 0 200px;
}
.ac-sce-node_body:not(.readonly) > div:nth-child(2) {
  flex: 0 0 140px;
}
.ac-sce-node_body:not(.readonly) > div:nth-child(3) {
  flex: 0 0 46px;
}
.ac-sce-node_body:not(.readonly) > div:nth-child(4) {
  flex: 0 0 260px;
}
.ac-sce-node_body:not(.readonly) > div:nth-child(5) {
  flex: 0 0 260px;
}
.ac-sce-node_body:not(.readonly) > div:nth-child(6) {
  flex: 0 0 48px;
}
.ac-sce-node_body .el-input__wrapper {
  --el-input-border-color: transparent;
  --el-input-bg-color: transparent;
  --el-disabled-border-color: transparent;
  --el-disabled-bg-color: transparent;
  padding: 1px 8px;
}

.ac-sce-ref > .ac-sce-node_children .ac-sce-node_drag {
  display: none;
}

.ac-sce-const > .ac-sce-node_drag {
  display: none;
}

.ac-sce-node .ac-selected {
  color: var(--el-color-primary);
  width: 120px;
}
</style>
