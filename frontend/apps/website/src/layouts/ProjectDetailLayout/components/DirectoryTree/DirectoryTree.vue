<template>
  <ToggleHeading :title="$t('app.interface.title')">
    <template #extra>
      <el-icon v-if="!isReader" class="cursor-pointer text-zinc-500" @click="onPopoverRefIconClick"><ac-icon-ep-plus /></el-icon>
    </template>
    <div ref="dir" :class="[ns.b(), { [ns.is('loading')]: isLoading }]" v-loading="isLoading">
      <ac-tree
        :data="apiDocTree"
        class="bg-transparent"
        node-key="id"
        empty-text=""
        ref="treeIns"
        :expand-on-click-node="false"
        :props="treeOptions"
        :draggable="!isReader"
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
            <div class="ac-tree-node__more" :class="{ active: data.id === activeNodeInfo?.id }" v-if="!isReader">
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

  <AIGenerateDocumentModal ref="aiPromptModalRef" @ok="onGenerateDocumenSuccess" />
</template>

<script setup lang="ts">
import documentIcon from '@/assets/images/doc-http@2x.png'
import AcTree from '@/components/AcTree'
import { useDocumentTree } from './useDocumentTree'
import { useDocumentPopoverMenu, PopoverMoreMenuType } from './useDocumentPopoverMenu'
import AIGenerateDocumentModal from '../AIGenerateDocumentModal.vue'
import { useActiveTree } from './useActiveTree'
import { useNamespace } from '@/hooks'
import { storeToRefs } from 'pinia'
import uesProjectStore from '@/store/project'

const schemaTree = inject('schemaTree') as any
const ns = useNamespace('catalog-tree')
const { isReader } = storeToRefs(uesProjectStore())

const {
  isLoading,
  treeIns,
  treeOptions,
  apiDocTree,
  handleTreeNodeClick,
  allowDrop,
  onMoveNode,
  onMoveNodeStart,
  updateTitle,
  initDocumentTree,
  redirecToDocumentEditPage,
  redirecToDocumentDetailPage,
} = useDocumentTree()

const aiPromptModalRef = ref<InstanceType<typeof AIGenerateDocumentModal>>()
const onGenerateDocumenSuccess = async (doc_id: any) => {
  schemaTree && (await schemaTree.reload())
  redirecToDocumentEditPage(doc_id)
}

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
} = useDocumentPopoverMenu(treeIns as any, aiPromptModalRef as any)

const { activeNode, reactiveNode } = useActiveTree(treeIns as any)

defineExpose({
  updateTitle,
  createNodeByData,
  reload: initDocumentTree,
  activeNode,
  reactiveNode,
  redirecToDocumentDetailPage,
})
</script>
