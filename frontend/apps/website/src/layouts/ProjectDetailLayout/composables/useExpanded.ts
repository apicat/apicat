import type Node from '@/components/AcTree/model/node'

export function useExpanded(_defaultExpandedKeys?: []) {
  const expandedKeysSet = ref(new Set<string | number>(_defaultExpandedKeys))
  const defaultExpandedKeys = computed(() => Array.from(expandedKeysSet.value))

  function handleNodeCollapse(node: Node) {
    expandedKeysSet.value.delete(node.data?.id)
  }

  function handleNodeExpand(node: Node) {
    expandedKeysSet.value.add(node.data?.id)
  }

  onBeforeUnmount(() => {
    expandedKeysSet.value.clear()
  })

  return {
    expandedKeysSet,
    defaultExpandedKeys,
    handleNodeCollapse,
    handleNodeExpand,
  }
}
