<script setup lang="ts">
import { useNamespace } from '@apicat/hooks'
import { ElTree } from 'element-plus'
import type { TreeData } from 'element-plus/es/components/tree/src/tree.type'
import type Node from '@/components/AcTree/model/node'

const props = withDefaults(
  defineProps<{
    historyRecord: HistoryRecord.TreeNode[]
    isLeaf?: (data: any) => boolean
  }>(),
  {
    historyRecord: () => [],
    isLeaf: (data: any) => !data.items,
  },
)

const emits = defineEmits(['onNodeClick'])

const treeClass = useNamespace('tree')
const nodeClass = useNamespace('tree-node')
const nodeContentClass = useNamespace('tree-content')
const activeHistoryID = ref<number>()
const treeOptions = { children: 'items', label: 'title', isLeaf: props.isLeaf }
const treeIns = ref<InstanceType<typeof ElTree>>()

function onTreeNodeClick(node: Node) {
  if (!node.isLeaf) {
    node.expanded = !node.expanded
    return
  }
  setCurrentKey(node.key)
  emits('onNodeClick', node)
}

function setCurrentKey(key: number) {
  activeHistoryID.value = key
  treeIns.value?.setCurrentKey(key)
}

defineExpose({
  setCurrentKey,
})
</script>

<template>
  <div class="flex flex-col h-full">
    <div class="flex-auto overflow-x-scroll scroll-content" :class="treeClass.b()">
      <ElTree
        ref="treeIns"
        :data="historyRecord"
        node-key="id"
        empty-text=""
        :props="treeOptions"
        :expand-on-click-node="false">
        <template #default="{ node }: { node: Node; data: TreeData }">
          <div
            :class="[nodeClass.b(), nodeClass.is('active', activeHistoryID === node.key)]"
            @click="onTreeNodeClick(node)">
            <div :class="nodeClass.e('main')">
              <div :id="`history_tree_node_${node.key}`" :class="nodeContentClass.b()">
                <i v-if="node.isLeaf" :class="nodeContentClass.e('icon')" class="ac-doc ac-iconfont" />
                <span :class="nodeContentClass.e('label')">{{ node.label }}</span>
              </div>
            </div>
          </div>
        </template>
      </ElTree>
    </div>
  </div>
</template>
