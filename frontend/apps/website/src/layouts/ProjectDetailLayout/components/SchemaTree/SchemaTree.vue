<template>
  <ToggleHeading title="模型">
    <template #extra>
      <el-icon class="cursor-pointer text-zinc-500" @click="onCreateSchemaMenuClick"><ac-icon-ep-plus /></el-icon>
    </template>
    <div ref="dir">
      <ac-tree
        :data="definitions"
        node-key="id"
        empty-text=""
        ref="treeIns"
        :expand-on-click-node="false"
        :props="treeOptions"
        draggable
        :allow-drop="allowDrop"
        @node-drag-start="onMoveNodeStart"
        @node-drop="onMoveNode"
      >
        <template #default="{ node, data }">
          <div class="flex justify-between ac-tree-node" :class="{ 'is-editable': data._extend.isEditable }">
            <div class="ac-tree-node__main" @click="handleTreeNodeClick(node, data, $event)">
              <div class="ac-doc-node" :class="{ 'is-active': data._extend.isCurrent }" :id="'schema_tree_node_' + data.id">
                <el-icon v-if="data._extend.isLeaf" class="ac-doc-node__icon" :size="17"><ac-icon-carbon-model-alt /></el-icon>
                <span class="ac-doc-node__label" v-show="!data._extend.isEditable" :title="data.name">{{ data.name }}</span>
              </div>
            </div>
            <div class="ac-tree-node__more" :class="{ active: data.id === activeNodeInfo?.id }">
              <el-icon v-show="!data._extend.isLeaf"><ac-icon-ep-plus /></el-icon>
              <span class="mx-1"></span>
              <el-icon @click="onPopoverRefIconClick($event, node)"><ac-icon-ep-more-filled /></el-icon>
            </div>
          </div>
        </template>
      </ac-tree>
    </div>
  </ToggleHeading>

  <el-popover :virtual-ref="popoverRefEl" trigger="click" virtual-triggering :visible="isShowPopoverMenu" width="auto">
    <PopperMenu :menus="popoverMenus" size="small" class="clear-popover-space" />
  </el-popover>
</template>

<script setup lang="ts">
import AcTree from '@/components/AcTree'
import { useSchemaPopoverMenu } from './useSchemaPopoverMenu'
import { useSchemaTree } from './useSchemaTree'
import { useActiveTree } from './useActiveTree'

const { treeIns, treeOptions, definitions, handleTreeNodeClick, allowDrop, onMoveNode, onMoveNodeStart, updateTitle } = useSchemaTree()
const { popoverMenus, popoverRefEl, isShowPopoverMenu, activeNodeInfo, onPopoverRefIconClick, onCreateSchemaMenuClick } = useSchemaPopoverMenu(treeIns as any)

const { activeNode } = useActiveTree(treeIns as any)

defineExpose({
  updateTitle,
  activeNode,
})
</script>
