import type { CollectionNode } from '@/typings/project'
import { storeToRefs } from 'pinia'
import { memoize } from 'lodash-es'

import AcTree from '@/components/AcTree'
import useDocumentStore from '@/store/document'
import { traverseTree } from '@apicat/shared'
import { TreeOptionProps } from '@/components/AcTree/tree.type'
import { DocumentTypeEnum } from '@/commons/constant'
import { useGoPage } from '@/hooks/useGoPage'
import { useActiveTree } from './useActiveTree'
import { DOCUMENT_DETAIL_NAME, DOCUMENT_EDIT_NAME } from '@/router'
import { useProjectId } from '@/hooks/useProjectId'
import { moveCollection } from '@/api/collection'

/**
 * 获取节点树最大深度
 */
const getTreeMaxDepth = memoize(function (node) {
  let maxLevel = 0
  traverseTree(
    (item: any) => {
      if (!item._extend.isLeaf) {
        maxLevel++
      }
    },
    [node] as CollectionNode[],
    { subKey: 'sub_nodes' }
  )
  return maxLevel
})

export const useDocumentTree = () => {
  const documentStore = useDocumentStore()
  const project_id = useProjectId()
  const { goDocumentDetailPage } = useGoPage()
  const route = useRoute()

  const { params } = route
  const { getApiDocTree } = documentStore
  const { apiDocTree } = storeToRefs(documentStore)

  const treeOptions: TreeOptionProps = {
    children: 'sub_nodes',
    label: 'title',
    class: (data): string => [(data as CollectionNode)._extend?.isLeaf ? 'is-doc' : 'is-dir'].join(' '),
    isLeaf: (data): boolean => (data as CollectionNode).type === DocumentTypeEnum.DOC,
  }

  const treeIns = ref<InstanceType<typeof AcTree>>()

  const { reactiveNode, activeNode } = useActiveTree(treeIns as any)

  /**
   * 目录树 点击
   */
  const handleTreeNodeClick = (node: any, source: any, e: Event) => {
    // 重命名输入框
    if ((e?.target as HTMLElement).tagName === 'INPUT') {
      e.preventDefault()
      return
    }

    // 文档点击
    if (source._extend.isLeaf) {
      goDocumentDetailPage(source.id)
      // router.push(getDocumentDetailPath(project_id as string, source.id))
      return
    }

    // 目录点击
    node.expanded = !node.expanded
  }

  /**
   * 允许拖拽
   */
  const allowDrop = (draggingNode: any, dropNode: any, type: any) => {
    const { data: dropNodeData } = dropNode
    const { data: draggingNodeData } = draggingNode

    // 不允许拖放在文件中
    if (dropNodeData._extend.isLeaf && type === 'inner') {
      return false
    }

    // 拖动目录时
    if (!draggingNodeData._extend.isLeaf && !dropNodeData._extend.isLeaf) {
      return getTreeMaxDepth(draggingNodeData) + dropNode.level <= 5
    }

    return true
  }

  let oldDraggingNodeInfo: any = null

  // 开始拖拽，记录旧节点位置数据
  const onMoveNodeStart = (draggingNode: any) => {
    const oldParent = draggingNode.parent

    oldDraggingNodeInfo = {
      oldPid: oldParent.id === 0 ? null : oldParent.key,
      oldChildIds: oldParent.childNodes.filter((item: any) => item.id !== draggingNode.id).map((item: any) => item.key),
    }
  }

  // 拖拽完成，更新新节点位置数据
  const onMoveNode = (draggingNode: any, dropNode: any, dropType: string) => {
    if (!oldDraggingNodeInfo) {
      return
    }

    const { oldPid, oldChildIds } = oldDraggingNodeInfo

    const isSeamLevel = oldPid === dropNode.parent.id && dropType !== 'inner'
    const newParent = treeIns.value?.getNode(draggingNode.data).parent

    // 手动展开父节点
    if (!newParent.expanded) {
      newParent.expanded = true
    }

    const newPid = newParent.id === 0 ? null : newParent.key
    const newChildIds = newParent.childNodes.map((item: any) => item.key)

    const sortParams = {
      target: {
        pid: newPid,
        ids: newChildIds,
      },
      origin: {
        pid: isSeamLevel ? newPid : oldPid,
        ids: isSeamLevel ? newChildIds : oldChildIds,
      },
    }

    oldDraggingNodeInfo = null
    moveCollection(project_id as string, sortParams)
  }

  const updateTitle = (id: any, title: string) => {
    const node = treeIns.value?.getNode(id)
    if (node && node?.data?.title) {
      node.data.title = title || node.data.title
    }
  }

  const initDocumentTree = async (activeDocId?: any) => {
    await getApiDocTree(project_id as string)
    if (route.name === DOCUMENT_DETAIL_NAME || route.name === DOCUMENT_EDIT_NAME) {
      params.doc_id ? activeNode(activeDocId || params.doc_id) : reactiveNode()
    }
  }

  const redirecToDocumentDetail = (activeId: any) => {
    goDocumentDetailPage(activeId)
    initDocumentTree(activeId)
  }

  onMounted(async () => initDocumentTree())

  return {
    treeIns,
    treeOptions,
    apiDocTree,

    handleTreeNodeClick,
    allowDrop,
    onMoveNodeStart,
    onMoveNode,
    updateTitle,

    initDocumentTree,

    redirecToDocumentDetail,
  }
}
