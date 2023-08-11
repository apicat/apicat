import { getCollectionList } from '@/api/collection'
import { DocumentTypeEnum } from '@/commons'
import { arrayToTree, flattenDeep, traverseForMarkSort } from '@/commons/helper'
import useApi from '@/hooks/useApi'
import { CollectionNode, EmptyStruct, Iteration } from '@/typings'

export const defaultProps = {
  label: 'title',
  children: 'items',
}

export const useIterationPlan = (iterationInfo: Ref<EmptyStruct<Iteration>>) => {
  const transferTreeRef = ref()
  const projectIdRef = ref<number | string | null>(null)
  const fromData = ref<CollectionNode[]>([])
  const toData = ref<CollectionNode[]>([])
  const [isLoadingForTree, execute] = useApi(getCollectionList)

  const onTransferTreeChange = (f: any, d: any) => {
    const from = arrayToTree(flattenDeep(f || [], defaultProps.children))
    const to = arrayToTree(flattenDeep(d || [], defaultProps.children))

    sortTree(from)
    sortTree(to)

    fromData.value = from
    toData.value = to

    // update collection_ids
    iterationInfo.value.collection_ids = transferTreeRef.value?.getFlattenValues().map((item: any) => item.id)
  }

  watch(projectIdRef, async (project_id) => {
    if (!project_id) {
      fromData.value = []
      toData.value = []
      return
    }

    const { id: iteration_id } = iterationInfo.value
    const allData = await execute(project_id, { iteration_id })
    const { from, to } = convertTransferTreeData(allData)
    fromData.value = from
    toData.value = to
  })

  return {
    projectIdRef,

    defaultProps,
    transferTreeRef,
    fromData,
    toData,

    isLoadingForTree,

    onTransferTreeChange,
  }
}

const convertTransferTreeData = (allData: any) => {
  traverseForMarkSort({ items: allData })

  const arr = flattenDeep(allData || [], defaultProps.children)

  const toDataArr: any = []
  const removeIndexSet = new Set()

  arr.forEach((item: any) => {
    if (item.selected) {
      item.items = []
      toDataArr.push(item)
      // 移除所有文档类型
      if (item.type !== DocumentTypeEnum.DIR) {
        removeIndexSet.add(item.id)
      }
    }
  })

  let fromDataArr = arr.filter((item: any) => !removeIndexSet.has(item.id))

  // 清空所有子集，重新组成新树
  fromDataArr.map((item: any) => {
    item.items = []
    return item
  })

  fromDataArr = arrayToTree(fromDataArr)

  // 移除所有空目录
  removeEmptyDirectory(fromDataArr, null, null)

  toDataArr.map((item: any) => {
    item.items = []
    return item
  })

  return {
    from: fromDataArr,
    to: toDataArr,
  }
}

const removeEmptyDirectory = (arr: any, node: any, parent: any) => {
  let i = arr.length

  while (i--) {
    const item = arr[i]
    removeEmptyDirectory(item.items, item, arr)
  }

  if (node && node.type === DocumentTypeEnum.DIR) {
    const child = node.items || []
    const len = child.length
    if (!len) {
      const idx = parent.indexOf(node)
      idx != -1 && parent.splice(idx, 1)
    }
  }
}

const sortTree = (nodes: any) => {
  ;(nodes || []).forEach((item: any) => {
    if (item.items && item.items.length) {
      sortTree(item.items)
    }
  })
  nodes.sort((pre: any, next: any) => pre.sort - next.sort)
}
