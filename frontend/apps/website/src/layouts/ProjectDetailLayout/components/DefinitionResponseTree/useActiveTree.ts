import { CollectionNode } from '@/typings/project'
import AcTree from '@/components/AcTree'
import useDefinitionResponseStore from '@/store/definitionResponse'
import { traverseTree } from '@apicat/shared'
import { storeToRefs } from 'pinia'
import scrollIntoView from 'smooth-scroll-into-view-if-needed'
import { useGoPage } from '@/hooks/useGoPage'

export const useActiveTree = (treeIns: Ref<InstanceType<typeof AcTree>>) => {
  const definitionResponseStore = useDefinitionResponseStore()
  const route = useRoute()
  const { goResponseDetailPage, goDocumentDetailPage } = useGoPage()
  const { responses } = storeToRefs(definitionResponseStore)
  const { params } = route
  const directoryTree = inject('directoryTree') as any

  // 启动切换文档选中
  watch(
    () => route.params.response_id,
    () => {
      activeNode(route.params.response_id)
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
    }, responses.value as any)

    const id = parseInt(nodeId as string, 10)
    const node = treeIns.value?.getNode(id)

    if (node && node.data) {
      ;(node.data as CollectionNode)._extend!.isCurrent = true
      treeIns.value?.setCurrentKey(id)
      const el = document.querySelector('#response_tree_node_' + id)
      el && scrollIntoView(el, { scrollMode: 'if-needed' })
    }
  }

  const reactiveNode = () => {
    if (!treeIns.value || !String(route.name).includes('definition.response')) {
      return
    }

    let hasCurrent = false
    traverseTree((item: CollectionNode) => {
      if (item._extend!.isCurrent) {
        hasCurrent = true
        return false
      }
    }, responses.value as any)

    // 没有选中模型时，进行自动切换
    if (!hasCurrent) {
      let node: any = null
      traverseTree((item: any) => {
        let _node = treeIns.value.getNode(item.id)

        if (_node && _node.data && _node.data._extend.isLeaf) {
          node = _node
          return false
        }
      }, responses.value)

      // 存在模型
      if (node) {
        params.response_id = node.key
        activeNode(node.key)
        // router.replace({ name: RESPONSE_DETAIL_NAME, params })
        goResponseDetailPage(node.key, true)
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
