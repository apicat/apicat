import { CollectionNode } from '@/typings/project'
import AcTree from '@/components/AcTree'
import useDefinitionStore from '@/store/definitionSchema'
import { traverseTree } from '@apicat/shared'
import { storeToRefs } from 'pinia'
import scrollIntoView from 'smooth-scroll-into-view-if-needed'
import { useGoPage } from '@/hooks/useGoPage'

export const useActiveTree = (treeIns: Ref<InstanceType<typeof AcTree>>) => {
  const definitionStore = useDefinitionStore()
  const route = useRoute()
  const { goSchemaDetailPage, goDocumentDetailPage } = useGoPage()
  const { definitions } = storeToRefs(definitionStore)
  const directoryTree = inject('directoryTree') as any

  // 启动切换文档选中
  watch(
    () => route.params.schema_id,
    () => {
      activeNode(route.params.schema_id)
    }
  )

  // 文档选中切换
  const activeNode = (nodeId?: any) => {
    // 清除选中
    traverseTree((item: CollectionNode) => {
      if (item._extend!.isCurrent) {
        item._extend!.isCurrent = false
        return false
      }
    }, definitions.value as any)

    const id = parseInt(nodeId as string, 10)
    const node = treeIns.value?.getNode(id)

    if (node && node.data) {
      ;(node.data as CollectionNode)._extend!.isCurrent = true
      treeIns.value?.setCurrentKey(id)
      // scrollIntoView
      const el = document.querySelector('#schema_tree_node_' + id)
      el && scrollIntoView(el, { scrollMode: 'if-needed' })
    }
  }

  const reactiveNode = () => {
    if (!treeIns.value || !String(route.name).includes('definition.schema')) {
      return
    }

    let hasCurrent = false
    traverseTree((item: CollectionNode) => {
      if (item._extend!.isCurrent) {
        hasCurrent = true
        return false
      }
    }, definitions.value as any)

    // 没有选中模型时，进行自动切换
    if (!hasCurrent) {
      let node: any = null
      traverseTree((item: any) => {
        let _node = treeIns.value.getNode(item.id)

        if (_node && _node.data && _node.data._extend.isLeaf) {
          node = _node
          return false
        }
      }, definitions.value)

      // 存在模型
      if (node) {
        // params.schema_id = node.key
        activeNode(node.key)
        // router.replace({ name: SCHEMA_DETAIL_NAME, params })
        goSchemaDetailPage(node.key, true)
      } else {
        // router.replace({ name: DOCUMENT_DETAIL_NAME, params: { project_id } })
        goDocumentDetailPage(undefined, true)
        setTimeout(() => directoryTree.reactiveNode && directoryTree.reactiveNode(), 0)
      }
    }
  }

  return {
    activeNode,
    reactiveNode,
  }
}
