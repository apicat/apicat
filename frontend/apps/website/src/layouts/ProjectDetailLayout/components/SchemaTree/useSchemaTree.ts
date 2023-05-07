import { useApi } from '@/hooks/useApi'
import type { CollectionNode } from '@/typings/project'
import { storeToRefs } from 'pinia'
import { memoize } from 'lodash-es'

import AcTree from '@/components/AcTree'
import { moveDefinition } from '@/api/definition'
import { traverseTree } from '@apicat/shared'
import { TreeOptionProps } from '@/components/AcTree/tree.type'
import { DocumentTypeEnum } from '@/commons/constant'
import { useActiveTree } from './useActiveTree'
import { useGoPage } from '@/hooks/useGoPage'
import { SCHEMA_DETAIL_NAME, SCHEMA_EDIT_NAME } from '@/router'
import useDefinitionStore from '@/store/definition'
import { useProjectId } from '@/hooks/useProjectId'

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
    { subKey: 'items' }
  )
  return maxLevel
})

/**
 * 此处逻辑和文档树逻辑可以进行优化
 * @returns
 */
export const useSchemaTree = () => {
  const definitionStore = useDefinitionStore()
  const project_id = useProjectId()
  const { goSchemaDetailPage, goSchemaEditPage } = useGoPage()
  const route = useRoute()
  const router = useRouter()
  const { params } = route
  const { getDefinitions } = definitionStore
  const { definitions } = storeToRefs(definitionStore)
  const [isLoading, getDefinitionsApi] = useApi(getDefinitions)()

  const treeOptions: TreeOptionProps = {
    children: 'items',
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
      goSchemaDetailPage(source.id)
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
    const newPid = newParent.id === null ? 0 : newParent.key
    const newChildIds = newParent.childNodes.map((item: any) => item.key)

    // 手动展开父节点
    if (!newParent.expanded) {
      newParent.expanded = true
    }

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
    moveDefinition(project_id as string, sortParams)
  }

  const updateTitle = (id: any, name: string) => {
    const node = treeIns.value?.getNode(id)
    if (node && node?.data?.name) {
      node.data.name = name || node.data.name
    }
  }

  const initSchemaTree = async (activeId?: any) => {
    await getDefinitionsApi(project_id as string)
    if (route.name === SCHEMA_DETAIL_NAME || route.name === SCHEMA_EDIT_NAME) {
      router.currentRoute.value.params.shcema_id ? activeNode(activeId || params.shcema_id) : reactiveNode()
    }
  }

  const redirecToSchemaEdit = (activeId: any) => {
    goSchemaEditPage(activeId)
    initSchemaTree(activeId)
  }

  onMounted(async () => await initSchemaTree())
  onUnmounted(() => definitionStore.$reset())
  return {
    isLoading,
    treeIns,
    treeOptions,
    definitions,

    handleTreeNodeClick,
    allowDrop,
    onMoveNodeStart,
    onMoveNode,
    updateTitle,

    initSchemaTree,

    redirecToSchemaEdit,
  }
}
