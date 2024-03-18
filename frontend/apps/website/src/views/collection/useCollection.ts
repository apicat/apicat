import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { storeToRefs } from 'pinia'
import type { PageModeCtx } from '../composables/usePageMode'
import useApi from '@/hooks/useApi'
import { useCollectionsStore } from '@/store/collections'
import { useProjectLayoutContext } from '@/layouts/ProjectDetailLayout'
import { getCollectionHistoryPath } from '@/router/history'
import { injectPagesMode } from '@/layouts/ProjectDetailLayout/composables/usePagesMode'
import { useTitleInputFocus } from '@/hooks/useTitleInputFocus'

let oldTitle = ''

export function useCollection(props: { project_id: string; collectionID: string }) {
  const { t } = useI18n()
  const router = useRouter()
  const { toggleMode, readonly, switchToWriteMode } = injectPagesMode('collection') as PageModeCtx
  const layoutContext = useProjectLayoutContext()
  const collectionIDRef = toRef(props, 'collectionID')
  const collectionStore = useCollectionsStore()
  const { collectionDetail: collection, loading: isLoading } = storeToRefs(collectionStore)
  const [isSaving, updateCollection, isSaveError] = useApi(collectionStore.updateCollection)
  const { inputRef, focus } = useTitleInputFocus()

  function handleTitleBlur() {
    const title = collection.value?.title || ''
    if (!title || !title.trim()) {
      collection.value!.title = oldTitle
      oldTitle = ''
    }
  }

  function handleContentUpdate(content: Array<any>) {
    if (collection.value) collection.value.content = content
  }

  function onShareCollectionBtnClick() {
    layoutContext.handleShareDocument!(props.project_id, props.collectionID)
  }

  function onExportCollectionBtnClick() {
    layoutContext.handleExportDocument!(props.project_id, props.collectionID)
  }

  function goCollectionhistory() {
    router.push({
      path: getCollectionHistoryPath(props.project_id, Number.parseInt(props.collectionID)),
      query: { ...router.currentRoute.value.query, iterationID: router.currentRoute.value.params.iterationID },
    })
  }

  watchDebounced(
    [() => collection.value?.id, () => collection.value?.title, () => collection.value?.content],
    async ([nId], [oId]) => {
      const n = collection.value

      if (readonly.value || !n) return

      // 还原旧的title时，不需要请求接口
      if (!oldTitle) {
        oldTitle = n?.title || ''
        return
      }

      if (nId === oId) {
        // title is empty
        if (!n.title || !n.title.trim()) {
          ElMessage.error(t('app.project.collection.page.edit.titleNull'))
          return
        }

        // backup old title
        oldTitle = n.title

        try {
          await updateCollection(props.project_id, n)
        } catch (error) {
          //
        }
      }
    },
    {
      debounce: 200,
    },
  )

  watch(
    collectionIDRef,
    async (id, oID) => {
      if (id === oID) return

      oldTitle = ''
      const collectionID = Number.parseInt(id)
      if (!Number.isNaN(collectionID)) {
        await collectionStore.getCollectionDetail(props.project_id, collectionID)
        oldTitle = collection.value?.title || ''
        if (!readonly.value) focus()
      }
    },
    {
      immediate: true,
    },
  )

  onBeforeUnmount(() => {
    collectionStore.collectionDetail = null
  })

  return {
    collection,
    isLoading,
    isSaving,
    isSaveError,

    handleTitleBlur,
    handleContentUpdate,
    onShareCollectionBtnClick,
    onExportCollectionBtnClick,
    goCollectionhistory,

    readonly,
    toggleMode,
    switchToWriteMode,
    titleInputRef: inputRef,
  }
}
