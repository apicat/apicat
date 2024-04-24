import { storeToRefs } from 'pinia'
import { traverseTree } from '@apicat/shared'
import type { ToggleHeading } from '@apicat/components'
import type AcTreeWrapper from '../AcTreeWrapper'
import { CurrentNodeContextKey, getParentNodeKeys } from '../../constants'
import { injectPagesMode } from '../../composables/usePagesMode'
import { useExpanded } from '../../composables/useExpanded'
import { useGoPage } from '@/hooks/useGoPage'
import useProjectStore from '@/store/project'
import type { TreeNodeData } from '@/components/AcTree/tree.type'
import { CollectionTypeEnum } from '@/commons'
import { useCollectionsStore } from '@/store/collections'
import { useParams } from '@/hooks/useParams'
import type { PageModeCtx } from '@/views/composables/usePageMode'
import type Node from '@/components/AcTree/model/node'

export function useSelectedNode(
  treeIns: Ref<InstanceType<typeof AcTreeWrapper> | undefined>,
  toggleHeadingRef: Ref<InstanceType<typeof ToggleHeading> | undefined>,
) {
  const { goCollectionPage } = useGoPage()
  const ctx = inject(CurrentNodeContextKey)
  const { switchToReadMode } = injectPagesMode('collection') as PageModeCtx
  const { projectID } = storeToRefs(useProjectStore())
  const { iterationID } = useParams()
  const { collections } = storeToRefs(useCollectionsStore())
  const { expandedKeysSet, ...rest } = useExpanded()

  async function selectedNodeWithGoPage(data: Node | TreeNodeData | undefined, isReplace = false) {
    // 跳转至默认页面
    if (!data) {
      await goCollectionPage(projectID.value!, '')
      return
    }

    const node = treeIns.value?.getNode(data)
    getParentNodeKeys(node).forEach(key => expandedKeysSet.value.add(key))
    if (node) {
      toggleHeadingRef.value?.expand()
      switchToReadMode()
      await goCollectionPage(projectID.value!, node.key, isReplace)
      ctx?.activeCollectionNode(node.key)
      treeIns.value?.setCurrentKey(node.key)
    }
  }

  function selectFirstNode() {
    let firstNode: CollectionAPI.ResponseCollection | undefined

    traverseTree<CollectionAPI.ResponseCollection>(
      (node) => {
        if (node.type !== CollectionTypeEnum.Dir) {
          if (iterationID.value && !node.selected) {
            return true
          }
          else {
            firstNode = node
            return false
          }
        }
        return true
      },
      collections.value,
      { subKey: 'items' },
    )

    return selectedNodeWithGoPage(firstNode, true)
  }

  function expandOnStartup() {
    if (ctx) {
      treeIns.value?.setCurrentKey(ctx.currentActiveNode.value.id)
      const node = treeIns.value?.getNode(ctx.currentActiveNode.value.id as any)
      getParentNodeKeys(node).forEach(key => expandedKeysSet.value.add(key))
    }
  }

  // 重新选中
  function reselectNode() {
    const currentID = ctx?.activeCollectionKey.value
    if (currentID !== -1 && !treeIns.value?.getNode(currentID as any))
      selectFirstNode()
  }

  return {
    selectedNodeWithGoPage,
    selectFirstNode,
    expandOnStartup,
    reselectNode,
    ...rest,
  }
}
