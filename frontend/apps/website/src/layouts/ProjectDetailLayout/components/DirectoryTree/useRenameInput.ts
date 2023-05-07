import { CollectionNode } from '@/typings/project'
import { debounce } from 'lodash-es'
import scrollIntoView from 'smooth-scroll-into-view-if-needed'
import { ActiveNodeInfo } from '@/typings/common'
import { updateCollection } from '@/api/collection'
import { useParams } from '@/hooks/useParams'

export const useRenameInput = (activeNodeInfo: Ref<ActiveNodeInfo | null>) => {
  const renameInputRef = ref<Nullable<HTMLInputElement>>(null)
  const { project_id } = useParams()
  /**
   * 重命名操作
   */
  const onRenameMenuClick = () => {
    const node = unref(activeNodeInfo)?.node
    const data = node?.data as CollectionNode
    if (data) {
      data._extend!.isEditable = true
      data._oldName = data.title
      renameInputFocus()
    }
  }

  const renameInputFocus = async () => {
    await nextTick()
    const inputEl = unref(renameInputRef)
    if (inputEl) {
      scrollIntoView(inputEl, { scrollMode: 'if-needed' })
      inputEl.focus()
      inputEl.setSelectionRange(0, inputEl.value.length)
    }
  }

  const onRenameInputEnterKeyUp = (e: Event) => e.target && (e.target as HTMLInputElement).blur()

  const onRenameInputBlur = debounce(async function (e: Event, source: CollectionNode) {
    source._extend!.isEditable = false

    // 进行数据还原
    if (!(e.target as HTMLInputElement).value && source._oldName) {
      source.title = source._oldName
      return
    }
    await updateCollection({ project_id, collection_id: source.id, title: source.title })
  }, 200)

  return {
    renameInputRef,

    renameInputFocus,
    onRenameInputEnterKeyUp,
    onRenameInputBlur,
    onRenameMenuClick,
  }
}
