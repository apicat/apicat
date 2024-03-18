<script setup lang="ts">
import { storeToRefs } from 'pinia'
import { ToggleHeading } from '@apicat/components'
import AcTreeWrapper from '../AcTreeWrapper'
import { CurrentNodeContextKey, PopoverMoreMenuType } from '../../constants'
import { useMenus } from './useMenus'
import { useSelectedNode } from './useSelectedNode'
import useProjectStore from '@/store/project'
import { useCollectionsStore } from '@/store/collections'
import type Node from '@/components/AcTree/model/node'

const ctx = inject(CurrentNodeContextKey)
const { isManager, isWriter, projectID } = storeToRefs(useProjectStore())
const collectionStore = useCollectionsStore()
const { collections } = storeToRefs(collectionStore)
const toggleHeadingRef = ref<InstanceType<typeof ToggleHeading>>()
const treeWrapper = ref<InstanceType<typeof AcTreeWrapper>>()
const { expandOnStartup, selectFirstNode, selectedNodeWithGoPage, defaultExpandedKeys, handleNodeCollapse, handleNodeExpand } = useSelectedNode(treeWrapper, toggleHeadingRef)

const { popoverRefEl, isShowPopoverMenu, popoverMenus, onPopoverIconClick, handleRenameCollection } = useMenus(
  treeWrapper,
  toggleHeadingRef,
)

function handleClickNode(node: Node) {
  selectedNodeWithGoPage(node)
}

async function handleSort(target: any, origin: any) {
  await collectionStore.sortCollections(projectID.value!, { target, origin })
}

defineExpose({
  selectFirstNode,
  expandOnStartup,
})
</script>

<template>
  <ToggleHeading ref="toggleHeadingRef" :title="$t('app.interface.title')">
    <template v-if="isManager || isWriter" #extra>
      <el-icon class="cursor-pointer text-zinc-500" @click="onPopoverIconClick">
        <ac-icon-ep-plus />
      </el-icon>
    </template>

    <AcTreeWrapper
      ref="treeWrapper"
      node-key="id"
      :props="{ class: (data) => (data.selected !== undefined && data.selected === false ? 'hidden' : '') }"
      :active-key="ctx?.activeCollectionKey.value"
      :datas="collections"
      :draggable="isManager || isWriter"
      :default-expanded-keys="defaultExpandedKeys"
      @node-collapse="handleNodeCollapse"
      @node-expand="handleNodeExpand"
      @rename="handleRenameCollection"
      @click="handleClickNode"
      @sort="handleSort"
    >
      <template #leafIcon>
        <i class="ac-doc ac-iconfont" />
      </template>

      <template v-if="isManager || isWriter" #moreMenu="{ node }">
        <el-icon v-show="!node.isLeaf" @click="onPopoverIconClick($event, node, PopoverMoreMenuType.ADD)">
          <ac-icon-ep-plus />
        </el-icon>
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
