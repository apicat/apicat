<script setup lang="ts">
import { ClickOutside as vClickOutside } from 'element-plus'
import { useNamespace } from '@apicat/hooks'
import { useRename } from './useRename'
import { useTree } from './useTree'
import AcTree from '@/components/AcTree'
import type Node from '@/components/AcTree/model/node'
import type { TreeData, TreeOptionProps } from '@/components/AcTree/tree.type'

export interface SortParams {
  parentID: number
  ids: number[]
}

export interface AcTreeWrapperProps {
  props?: TreeOptionProps
  datas: TreeData
  nodeKey: string
  activeKey?: string | number | null | undefined
  draggable?: boolean
  defaultExpandedKeys?: Array<string | number>
}

export interface AcTreeWrapperEmits {
  (e: 'click' | 'node-collapse' | 'node-expand', node: Node): void
  (e: 'rename', node: Node, newName: string, oldName: string): void
  (e: 'sort', target: SortParams, origin: SortParams): void
}

const props = defineProps<AcTreeWrapperProps>()
const emits = defineEmits<AcTreeWrapperEmits>()

const treeClass = useNamespace('tree')
const nodeClass = useNamespace('tree-node')
const nodeContentClass = useNamespace('tree-content')
const operationInfo = ref<{ node?: Node; key?: any }>({})

const {
  NODE_UUID_PREFIX,
  treeIns,
  treeWrapper,
  treeOptions,
  allowDrop,
  handleDragStart,
  handleDragEnd,
  handleTreeNodeClick,
  scrollIntoViewByKey,
  append,
  remove,
  insertAfter,
  insertBefore,
  getNode,
  setCurrentKey,
} = useTree(props, emits)

const { currentRenameNode, renameInputRef, nodeName, onRenameInputBlur, onRenameInputEnterKeyUp, onRenameNode } = useRename(treeIns, emits)

defineExpose({
  append,
  remove,
  insertAfter,
  insertBefore,
  getNode,
  setCurrentKey,
  rename: onRenameNode,
  scrollIntoViewByKey,
})

function onClickMoreAreaOutside() {
  operationInfo.value.key = undefined
}

function onClickMoreArea(node: Node) {
  operationInfo.value = { node, key: node.key }
}
</script>

<template>
  <div ref="treeWrapper" :class="treeClass.b()">
    <AcTree
      ref="treeIns"
      empty-text=""
      :data="datas"
      :node-key="nodeKey"
      :props="treeOptions"
      :expand-on-click-node="false"
      :auto-expand-parent="false"
      :draggable="draggable"
      :allow-drag="(node: Node) => node !== currentRenameNode"
      :allow-drop="allowDrop"
      :default-expanded-keys="defaultExpandedKeys"
      @node-expand="(data, node:Node) => emits('node-expand', node)"
      @node-collapse="(data, node:Node) => emits('node-collapse', node)"
      @node-drag-start="handleDragStart"
      @node-drop="handleDragEnd"
    >
      <template #default="{ node, data }: { node: Node; data: TreeData }">
        <div :class="[nodeClass.b(), nodeClass.is('active', node.isLeaf && activeKey === node.key)]">
          <div :class="nodeClass.e('main')" @click="handleTreeNodeClick(node, $event)">
            <div :id="`${NODE_UUID_PREFIX}${node.key}`" :class="nodeContentClass.b()">
              <span v-if="node.isLeaf" :class="nodeContentClass.e('icon')">
                <slot name="leafIcon" :node="node" :data="data" />
              </span>
              <span v-if="currentRenameNode !== node" :class="nodeContentClass.e('label')" :title="node.label">{{
                node.label
              }}</span>
              <el-input
                v-if="currentRenameNode === node"
                ref="renameInputRef"
                v-model="nodeName"
                size="small"
                maxlength="255"
                :class="nodeContentClass.e('input')"
                @keyup.enter="onRenameInputEnterKeyUp"
                @blur="onRenameInputBlur"
              />
            </div>
          </div>

          <div
            v-click-outside="onClickMoreAreaOutside"
            :class="[nodeClass.e('more'), nodeClass.is('active', operationInfo.key === node.key)]"
            @click="onClickMoreArea(node)"
          >
            <slot name="moreMenu" :node="node" :data="data" :operate-info="operationInfo" />
          </div>
        </div>
      </template>
    </AcTree>
  </div>
</template>
