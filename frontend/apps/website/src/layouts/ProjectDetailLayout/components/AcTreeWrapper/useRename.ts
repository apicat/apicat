import type { AcTreeWrapperEmits } from './AcTreeWrapper.vue'
import type AcTree from '@/components/AcTree'
import type Node from '@/components/AcTree/model/node'
import type { TreeKey, TreeNodeData } from '@/components/AcTree/tree.type'

export function useRename(treeIns: Ref<InstanceType<typeof AcTree> | undefined>, emits: AcTreeWrapperEmits) {
  const renameInputRef = ref<Nullable<HTMLInputElement>>(null)
  const currentRenameNode = ref<Node>()
  const nodeName = ref('')

  const renameInputFocus = async () => {
    await nextTick()
    setTimeout(() => {
      const inputEl = unref(renameInputRef)
      if (inputEl) {
        inputEl.focus()
        inputEl.select()
      }
    }, 0)
  }

  const onRenameInputBlur = async () => {
    const newName = (unref(nodeName) || '').trim()
    const node = unref(currentRenameNode)

    if (node && newName && newName !== node.label)
      emits('rename', node, newName, node.label)

    currentRenameNode.value = undefined
    nodeName.value = ''
  }

  const onRenameInputEnterKeyUp = (e: Event) => e.target && (e.target as HTMLInputElement).blur()

  /**
   * 重命名操作
   */
  const onRenameNode = (data: TreeKey | TreeNodeData) => {
    const node = treeIns.value?.getNode(data)
    currentRenameNode.value = node
    nodeName.value = node?.label || ''
    renameInputFocus()
  }

  return {
    nodeName,
    currentRenameNode,
    renameInputRef,

    renameInputFocus,
    onRenameInputEnterKeyUp,
    onRenameInputBlur,
    onRenameNode,
  }
}
