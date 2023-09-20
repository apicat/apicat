<template>
  <div class="flex flex-col h-full">
    <div class="flex-auto overflow-x-scroll scroll-content" ref="dir">
      <ac-tree
        :data="historyRecordTree"
        class="bg-transparent"
        node-key="id"
        empty-text=""
        ref="treeIns"
        :expand-on-click-node="false"
        :props="{ children: 'sub_nodes', label: 'title', class: customNodeClass, isLeaf: customNodeLeaf }"
      >
        <template #default="{ node, data }">
          <div class="el-tree-node__bg"></div>

          <div class="flex justify-between ac-tree-node">
            <div class="ac-tree-node__main" @click="handleTreeNodeClick(node, data, $event)">
              <div class="ac-doc-node" :class="{ 'is-active': data._extend.isCurrent }" :id="'history_schema_tree_node_' + data.id">
                <el-icon v-if="data._extend.isLeaf" class="ac-doc-node__icon" :size="17"><ac-icon-carbon-model-alt /></el-icon>
                <span class="ac-doc-node__label" :title="data.name">{{ data.name }}</span>
              </div>
            </div>
          </div>
        </template>
      </ac-tree>
    </div>
  </div>
</template>

<script lang="ts" setup>
import AcTree from '@/components/AcTree'
import scrollIntoView from 'smooth-scroll-into-view-if-needed'
import { useDefinitionSchemaStore } from '@/store/definitionSchema'
import { storeToRefs } from 'pinia'
import { traverseTree } from '@apicat/shared'
import { DocumentTypeEnum } from '@/commons/constant'
import { useParams } from '@/hooks/useParams'
import { CollectionNode } from '@/typings/project'

const $route = useRoute()
const $router = useRouter()
const { project_id, computedRouteParams } = useParams()
const { params } = $route

const definitionSchemaStore = useDefinitionSchemaStore()
const { historyRecordTree } = storeToRefs(definitionSchemaStore)

const treeIns: any = ref(null)
const dir: Ref<HTMLDivElement | null> = ref(null)

const customNodeClass = (data: any) => (data._extend.isLeaf ? 'is-doc' : 'is-dir')
const customNodeLeaf = (data: any) => data.type !== DocumentTypeEnum.DIR

const handleTreeNodeClick = (node: any, source: any, e: any) => {
  if (e.target.tagName === 'INPUT') {
    return
  }

  // 文档点击
  if (source._extend.isLeaf) {
    onDocumentClick(source)
    return
  }
  // 目录点击
  node.expanded = !node.expanded
}

const onDocumentClick = (source: any) => {
  const { history_id } = $router.currentRoute.value.params
  // 同一篇文档，且为详情页，不进行任何操作
  if (source.id === parseInt(history_id as string, 10)) {
    return
  }

  activeNode(source.id)

  $router.push({
    name: 'history.schema.detail',
    params: { ...params, history_id: source.id },
  })
}

// 文档选中切换
const activeNode = (nodeId?: any) => {
  // 清除选中
  traverseTree(
    (item: CollectionNode) => {
      if (item._extend!.isCurrent) {
        item._extend!.isCurrent = false
        return false
      }
    },
    historyRecordTree.value as any,
    { subKey: 'sub_nodes' }
  )

  const id = parseInt(nodeId as string, 10)
  const node = treeIns.value?.getNode(id)
  if (node && node.data) {
    ;(node.data as CollectionNode)._extend!.isCurrent = true
    treeIns.value?.setCurrentKey(id)
    const el = document.querySelector('#history_schema_tree_node_' + id)
    el && scrollIntoView(el, { scrollMode: 'if-needed' })
  }
}

const reactiveNode = () => {
  if (!treeIns.value) {
    return
  }

  let hasCurrent = false
  traverseTree(
    (item: any) => {
      if (item._extend!.isCurrent) {
        hasCurrent = true
        return false
      }
    },
    historyRecordTree.value as any,
    { subKey: 'sub_nodes' }
  )

  // 没有选中文档时，进行自动切换
  if (!hasCurrent) {
    let node: any = null
    traverseTree(
      (item: any) => {
        let _node = treeIns.value.getNode(item.id)

        if (_node && _node.data._extend!.isLeaf) {
          node = _node
          return false
        }
      },
      historyRecordTree.value,
      { subKey: 'sub_nodes' }
    )

    // 存在文档
    if (node) {
      params.history_id = node.key
      activeNode(node.key)
    }

    dir.value?.scrollTo(0, 0)
    $router.replace({ name: 'history.schema.detail', params })
  }
}

onMounted(async () => {
  const { schema_id } = unref(computedRouteParams)
  await definitionSchemaStore.getSchemaHistoryRecordList(project_id, schema_id)
  params.history_id ? activeNode(params.history_id) : reactiveNode()
})
</script>
