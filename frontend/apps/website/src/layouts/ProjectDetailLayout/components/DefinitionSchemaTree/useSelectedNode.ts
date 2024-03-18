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
import { SchemaTypeEnum } from '@/commons'
import useDefinitionSchemaStore from '@/store/definitionSchema'
import type { PageModeCtx } from '@/views/composables/usePageMode'

export function useSelectedNode(
  treeIns: Ref<InstanceType<typeof AcTreeWrapper> | undefined>,
  toggleHeadingRef: Ref<InstanceType<typeof ToggleHeading> | undefined>,
) {
  const { goSchemaPage } = useGoPage()
  const ctx = inject(CurrentNodeContextKey)
  const { projectID } = storeToRefs(useProjectStore())
  const { switchToReadMode } = injectPagesMode('schema') as PageModeCtx
  const { schemas } = storeToRefs(useDefinitionSchemaStore())
  const { expandedKeysSet, ...rest } = useExpanded()

  function selectedNodeWithGoPage(data: Node | TreeNodeData) {
    const node = treeIns.value?.getNode(data)
    getParentNodeKeys(node).forEach(key => expandedKeysSet.value.add(key))
    if (node) {
      toggleHeadingRef.value?.expand()
      switchToReadMode()
      goSchemaPage(projectID.value!, node.key, false).then(() => {
        ctx?.activeSchemaNode(node.key)
        treeIns.value?.setCurrentKey(node.key)
      })
    }
  }

  function selectFirstNode() {
    traverseTree<Definition.SchemaNode>(
      (node) => {
        if (node.type !== SchemaTypeEnum.Category) {
          selectedNodeWithGoPage(node)
          return false
        }
        return true
      },
      schemas.value,
      { subKey: 'items' },
    )
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
    const currentID = ctx?.activeSchemaKey.value
    if (currentID !== undefined && !treeIns.value?.getNode(currentID as any))
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
