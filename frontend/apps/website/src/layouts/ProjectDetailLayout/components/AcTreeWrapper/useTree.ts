import scrollIntoView from 'smooth-scroll-into-view-if-needed'
import { uuid } from '@apicat/shared'
import type { AcTreeWrapperEmits, AcTreeWrapperProps } from './AcTreeWrapper.vue'
import type AcTree from '@/components/AcTree'
import type { DropType, TreeData, TreeKey, TreeNodeData, TreeOptionProps } from '@/components/AcTree/tree.type'
import { CollectionTypeEnum, checkDepth } from '@/commons'
import Node from '@/components/AcTree/model/node'

// 获取节点树最大深度(此方法有问题，应该获取被计算node 子集最大level - 被计算node资深level)
// const getTreeMaxDepth = createTreeMaxDepthFn<Node>('childNodes')

export function useTree(props: AcTreeWrapperProps, emits: AcTreeWrapperEmits) {
  const treeIns = ref<InstanceType<typeof AcTree>>()
  const treeWrapper = ref<HTMLDivElement>()
  // 节点前缀
  const NODE_UUID_PREFIX = `tree_node_${uuid()}_`

  // 目录树基础配置
  const treeOptions = computed<TreeOptionProps>(() => ({
    children: 'items',
    label: 'title',
    isLeaf: (data): boolean => data.type !== CollectionTypeEnum.Dir,
    ...(props.props || {}),
  }))

  // 节点点击
  const handleTreeNodeClick = (node: Node, e: Event) => {
    // 重命名输入框
    if ((e?.target as HTMLElement).tagName === 'INPUT') {
      e.preventDefault()
      return
    }

    // 文档点击
    if (node.isLeaf) {
      emits('click', node)
      return
    }
    // 目录点击
    node.expanded = !node.expanded
    emits(node.expanded ? 'node-expand' : 'node-collapse', node)
  }

  // 允许拖拽
  const allowDrop = (draggingNode: Node, dropNode: Node, type: DropType) => {
    // 不允许拖放在文件中
    if (dropNode.isLeaf && type === 'inner')
      return false

    // 拖动目录时
    if (!draggingNode.isLeaf) {
      // return getTreeMaxDepth(draggingNode) + dropNode.level <= 5
      return checkDepth(draggingNode, 'childNodes', data => data.isLeaf || false) + dropNode.level <= 5
    }

    return true
  }

  let oldDraggingNodeInfo: any = null

  // 开始拖拽，记录旧节点位置数据
  const handleDragStart = (draggingNode: Node) => {
    const oldParent = draggingNode.parent
    oldDraggingNodeInfo = {
      oldPid: oldParent.id === 0 ? null : oldParent.key,
      oldChildIds: oldParent.childNodes.filter((item: any) => item.id !== draggingNode.id).map((item: any) => item.key),
    }
  }

  // 拖拽完成，更新新节点位置数据
  const handleDragEnd = (draggingNode: any, dropNode: any, dropType: string) => {
    if (!oldDraggingNodeInfo)
      return

    const { oldPid, oldChildIds } = oldDraggingNodeInfo

    const isSeamLevel = oldPid === dropNode.parent.id && dropType !== 'inner'
    const newParent = treeIns.value?.getNode(draggingNode.data).parent

    // 手动展开父节点
    if (!newParent.expanded)
      newParent.expanded = true

    const newPid = newParent.id === 0 ? null : newParent.key
    const newChildIds = newParent.childNodes.map((item: any) => item.key)

    emits(
      'sort',
      {
        parentID: newPid,
        ids: newChildIds,
      },
      {
        parentID: isSeamLevel ? newPid : oldPid,
        ids: isSeamLevel ? newChildIds : oldChildIds,
      },
    )

    oldDraggingNodeInfo = null
  }

  // 滚动到指定节点
  const scrollIntoViewByKey = (key: TreeKey) => {
    const node = treeIns.value?.getNode(key)
    nextTick(() => {
      const el = document.querySelector(`#${NODE_UUID_PREFIX}${node?.key}`)
      el && setTimeout(() => scrollIntoView(el, { scrollMode: 'if-needed' }), 0)
    })
  }

  // proxy tree method
  const insertBefore = (data: TreeNodeData, refData: TreeData | TreeKey) => {
    treeIns.value?.insertBefore(data, refData)
    const parentNode = treeIns.value?.getNode(refData).parent
    expandByNode(parentNode)

    scrollIntoViewByKey(data[props.nodeKey])
  }

  const insertAfter = (data: TreeNodeData, refData: TreeData | TreeKey) => {
    treeIns.value?.insertAfter(data, refData)
    const parentNode = treeIns.value?.getNode(refData).parent
    expandByNode(parentNode)
    scrollIntoViewByKey(data[props.nodeKey])
  }

  const append = (data: TreeNodeData, parentNode: Node | TreeNodeData | TreeKey) => {
    treeIns.value?.append(data, parentNode)
    expandByNode(parentNode)
    scrollIntoViewByKey(data[props.nodeKey])
  }

  const getNode = (data: TreeKey | TreeNodeData): Node | undefined => treeIns.value?.getNode(data)

  const setCurrentKey = (key?: TreeKey) => {
    if (key === undefined)
      return

    treeIns.value?.setCurrentKey(key)
    scrollIntoViewByKey(key)
  }

  const remove = (node: Node | TreeNodeData) => treeIns.value?.remove(node)

  const expandByNode = (node: any) => {
    if (node instanceof Node && node.parent) {
      emits('node-expand', node)
      node.expand()
    }
  }

  return {
    NODE_UUID_PREFIX,

    treeWrapper,
    treeIns,
    treeOptions,

    allowDrop,
    scrollIntoViewByKey,
    handleTreeNodeClick,
    handleDragStart,
    handleDragEnd,

    // proxy method
    append,
    insertBefore,
    insertAfter,
    remove,
    getNode,
    setCurrentKey,
  }
}
