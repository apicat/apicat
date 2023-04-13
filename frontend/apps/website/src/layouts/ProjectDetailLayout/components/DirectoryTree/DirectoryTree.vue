<template>
  <ToggleHeading title="接口">
    <template #extra>
      <el-icon class="cursor-pointer text-zinc-500" @click="onPopoverRefIconClick"><ac-icon-ep-plus /></el-icon>
    </template>
    <div ref="dir">
      <ac-tree
        :data="apiDocTree"
        class="bg-transparent"
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
          <div class="el-tree-node__bg"></div>
          <div class="flex justify-between ac-tree-node" :class="{ 'is-editable': data._extend.isEditable }">
            <div class="ac-tree-node__main" @click="handleTreeNodeClick(node, data, $event)">
              <div class="ac-doc-node" :class="{ 'is-active': data._extend.isCurrent }" :id="'tree_node_' + data.id">
                <img v-if="data._extend.isLeaf" class="ac-doc-node__icon" :src="documentIcon" />
                <span class="ac-doc-node__label" v-show="!data._extend.isEditable" :title="data.title">{{ data.title }}</span>
                <input
                  type="text"
                  ref="renameInputRef"
                  class="ac-doc-node__input el-input el-input__inner"
                  v-if="data._extend.isEditable"
                  v-model="data.title"
                  maxlength="255"
                  @keyup.enter="onRenameInputEnterKeyUp"
                  @blur="onRenameInputBlur($event, data)"
                />
              </div>
            </div>
            <div class="ac-tree-node__more" :class="{ active: data.id === activeNodeInfo?.id }">
              <el-icon v-show="!data._extend.isLeaf" @click="onPopoverRefIconClick($event, node, PopoverMoreMenuType.ADD)"><ac-icon-ep-plus /></el-icon>
              <span class="mx-1"></span>
              <el-icon @click="onPopoverRefIconClick($event, node, PopoverMoreMenuType.MORE)"><ac-icon-ep-more-filled /></el-icon>
            </div>
          </div>
        </template>
      </ac-tree>
    </div>
  </ToggleHeading>

  <el-popover :virtual-ref="popoverRefEl" trigger="click" virtual-triggering :visible="isShowPopoverMenu" width="auto">
    <PopperMenu :menus="popoverMenus" :size="popoverMenuSize" class="clear-popover-space" />
  </el-popover>
</template>

<script setup lang="ts">
import documentIcon from '@/assets/images/doc-http@2x.png'
import AcTree from '@/components/AcTree'
import { useDocumentTree } from './useDocumentTree'
import { useDocumentPopoverMenu, PopoverMoreMenuType } from './useDocumentPopoverMenu'

const { treeIns, treeOptions, apiDocTree, handleTreeNodeClick, allowDrop, onMoveNode, onMoveNodeStart, updateTitle } = useDocumentTree()

const {
  popoverMenus,
  popoverRefEl,
  isShowPopoverMenu,
  activeNodeInfo,
  popoverMenuSize,
  renameInputRef,
  onPopoverRefIconClick,
  onRenameInputEnterKeyUp,
  createNodeByData,
  onRenameInputBlur,
} = useDocumentPopoverMenu(treeIns as any)

defineExpose({
  updateTitle,
  createNodeByData,
})
</script>
