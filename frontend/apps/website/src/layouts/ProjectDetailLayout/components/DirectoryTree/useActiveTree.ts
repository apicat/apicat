import { CollectionNode } from '@/typings/project'
import AcTree from '@/components/AcTree'
import useDocumentStore from '@/store/document'
import { traverseTree } from '@apicat/shared'
import { storeToRefs } from 'pinia'
import scrollIntoView from 'smooth-scroll-into-view-if-needed'
import { useGoPage } from '@/hooks/useGoPage'
import { useIterationStore } from '@/store/iteration'

export const useActiveTree = (treeIns: Ref<InstanceType<typeof AcTree>>) => {
  const documentStore = useDocumentStore()
  const iterationStore = useIterationStore()
  const route = useRoute()
  const { goDocumentDetailPage } = useGoPage()
  const { apiDocTree } = storeToRefs(documentStore)

  // 启动切换文档选中
  watch(
    () => route.params.doc_id,
    () => activeNode(route.params.doc_id)
  )

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
      apiDocTree.value as any,
      { subKey: 'items' }
    )

    const activeNodeId = nodeId

    const id = parseInt(activeNodeId as string, 10)
    const node = treeIns.value?.getNode(id)

    if (node && node.data) {
      ;(node.data as CollectionNode)._extend!.isCurrent = true
      treeIns.value?.setCurrentKey(id)
    }
    // scrollIntoView
    const el = document.querySelector('#tree_node_' + id)
    el && scrollIntoView(el, { scrollMode: 'if-needed' })
  }

  const reactiveNode = () => {
    if (!treeIns.value || !String(route.name).includes('document.')) {
      return
    }
    let hasCurrent = false
    traverseTree(
      (item: CollectionNode) => {
        if (item._extend!.isCurrent) {
          hasCurrent = true
          return false
        }
      },
      apiDocTree.value as any,
      { subKey: 'items' }
    )

    // 没有选中文档时，进行自动切换
    if (!hasCurrent) {
      let node: any = null
      traverseTree(
        (item: any) => {
          let _node = treeIns.value.getNode(item.id)
          // todo 迭代
          if (iterationStore.isIterationRoute) {
            if (_node && _node.data && _node.data._extend.isLeaf && _node.data.selected) {
              node = _node
              return false
            }
          } else {
            if (_node && _node.data && _node.data._extend.isLeaf) {
              node = _node
              return false
            }
          }
        },
        apiDocTree.value,
        { subKey: 'items' }
      )

      // 存在文档
      if (node) {
        activeNode(node.key)
        goDocumentDetailPage(node.key, true)
        // router.replace({ name: DOCUMENT_DETAIL_NAME, params })
        return
      }

      goDocumentDetailPage(undefined, true)
      // router.replace(getProjectDetailPath(route.params.project_id as string))
    }
  }

  return {
    activeNode,
    reactiveNode,
  }
}
