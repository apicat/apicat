<script setup lang="ts">
import { storeToRefs } from 'pinia'
import { ToggleHeading } from '@apicat/components'
import AcTreeWrapper from '../AcTreeWrapper'
import { CurrentNodeContextKey, PopoverMoreMenuType } from '../../constants'
import { useMenus } from './useMenus'
import { useSelectedNode } from './useSelectedNode'
import useProjectStore from '@/store/project'
import useDefinitionSchemaStore from '@/store/definitionSchema'
import type Node from '@/components/AcTree/model/node'

const { isManager, isWriter, projectID } = storeToRefs(useProjectStore())
const definitionSchemaStore = useDefinitionSchemaStore()
const { schemas } = storeToRefs(definitionSchemaStore)
const toggleHeadingRef = ref()
const treeWrapper = ref<InstanceType<typeof AcTreeWrapper>>()
const { popoverRefEl, isShowPopoverMenu, popoverMenus, onPopoverIconClick, handleRename } = useMenus(
  treeWrapper,
  toggleHeadingRef,
)
const { expandOnStartup, selectFirstNode, selectedNodeWithGoPage, defaultExpandedKeys, handleNodeCollapse, handleNodeExpand } = useSelectedNode(treeWrapper, toggleHeadingRef)
const ctx = inject(CurrentNodeContextKey)

function handleClickNode(node: Node) {
  selectedNodeWithGoPage(node)
}

async function handleSort(target: any, origin: any) {
  await definitionSchemaStore.moveSchema(projectID.value!, { target, origin })
}

defineExpose({
  selectFirstNode,
  expandOnStartup,
})
</script>

<template>
  <ToggleHeading ref="toggleHeadingRef" :title="$t('app.schema.title')">
    <template v-if="isManager || isWriter" #extra>
      <el-icon class="cursor-pointer text-zinc-500" @click="onPopoverIconClick">
        <ac-icon-ep-plus />
      </el-icon>
    </template>

    <AcTreeWrapper
      ref="treeWrapper"
      node-key="id"
      :props="{ label: 'name' }"
      :active-key="ctx?.activeSchemaKey.value"
      :datas="schemas"
      :draggable="isManager || isWriter"
      :default-expanded-keys="defaultExpandedKeys"
      @node-collapse="handleNodeCollapse"
      @node-expand="handleNodeExpand"
      @rename="handleRename"
      @click="handleClickNode"
      @sort="handleSort"
    >
      <template #leafIcon>
        <Iconfont class="ac-doc-node__icon" icon="ac-model" />
      </template>
      <template v-if="isManager || isWriter" #moreMenu="{ node }">
        <el-icon v-show="!node.isLeaf" @click="onPopoverIconClick($event, node, PopoverMoreMenuType.ADD)">
          <ac-icon-ep-plus />
        </el-icon>
        <span class="mx-1" />
        <el-icon @click="onPopoverIconClick($event, node, PopoverMoreMenuType.MORE)">
          <ac-icon-ep-more-filled />
        </el-icon>
      </template>
    </AcTreeWrapper>
  </ToggleHeading>

  <el-popover
    width="auto"
    transition="fade-fast"
    trigger="click"
    virtual-triggering
    :virtual-ref="popoverRefEl"
    :visible="isShowPopoverMenu"
    :show-arrow="false"
  >
    <PopperMenu :menus="popoverMenus" class="normal-popover-space" size="thin" />
  </el-popover>
</template>
