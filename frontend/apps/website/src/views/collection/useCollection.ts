import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { storeToRefs } from 'pinia'
import type { PageModeCtx } from '../composables/usePageMode'
import { useAITips } from './useAITips'
import useApi from '@/hooks/useApi'
import { useCollectionsStore } from '@/store/collections'
import { useProjectLayoutContext } from '@/layouts/ProjectDetailLayout'
import { getCollectionHistoryPath } from '@/router/history'
import { injectPagesMode } from '@/layouts/ProjectDetailLayout/composables/usePagesMode'
import { useTitleInputFocus } from '@/hooks/useTitleInputFocus'
import { injectAsyncInitTask } from '@/hooks/useWaitAsyncTask'

let oldTitle = ''

export function useCollection(props: { project_id: string, collectionID: string }) {
  const { t } = useI18n()
  const router = useRouter()
  const { toggleMode, readonly, switchToWriteMode } = injectPagesMode('collection') as PageModeCtx
  const layoutContext = useProjectLayoutContext()
  const collectionIDRef = toRef(props, 'collectionID')
  const collectionStore = useCollectionsStore()
  const { collectionDetail: collection, loading: isLoading } = storeToRefs(collectionStore)
  const [isSaving, updateCollection, isSaveError] = useApi(collectionStore.updateCollection)
  const { inputRef, focus } = useTitleInputFocus()
  const { preCollection, docTitle, isAIMode, isShowAIStyle, isShowAIStyleForTitle, handleEditorEvent, handleTitleBlur: triggerTitleBlur, cancelAITips, onDocumentLayoutClick } = useAITips(props.project_id, collection, readonly, updateCollection, collectionIDRef)

  function handleTitleBlur() {
    const title = collection.value?.title || ''
    if (!title || !title.trim()) {
      collection.value!.title = oldTitle
      oldTitle = ''
    }
    triggerTitleBlur()
  }

  async function handleContentUpdate(content: Array<any>) {
    if (!collection.value)
      return
    collection.value.content = content
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

  watchDebounced([() => collection.value?.id, () => collection.value?.title, () => collection.value?.content], async ([nId], [oId]) => {
    const n = collection.value

    if (readonly.value || !n)
      return

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
      // sync docTitle
      docTitle.value = n.title
      if (!isAIMode.value)
        await updateCollection(props.project_id, n)
    }
  }, { debounce: 500 })

  const asyncTaskCtx = injectAsyncInitTask()
  let isImmidiate = true
  watch(
    collectionIDRef,
    async (id, oID) => {
      if (id === oID)
        return

      oldTitle = ''
      const collectionID = Number.parseInt(id)
      if (!Number.isNaN(collectionID)) {
        const task = collectionStore.getCollectionDetail(props.project_id, collectionID)
        if (isImmidiate) {
          isImmidiate = false
          asyncTaskCtx!.addTask(task)
        }
        // 等待数据加载完成，获取详情
        await task
        // 重置AI提示
        cancelAITips()
        if (collection.value) {
          preCollection.value = JSON.parse(JSON.stringify(collection.value))
          docTitle.value = ''
        }

        oldTitle = collection.value?.title || ''
        if (!readonly.value)
          focus()

        // // 如果内容为空，则可以开启AI推理模式,备份当前collection
        // isAIMode.value = isEmptyContent(collection.value?.content)
        // if (isAIMode.value)
        //   preCollection.value = JSON.parse(JSON.stringify(collection.value))
      }
    },
    {
      immediate: true,
    },
  )

  onBeforeUnmount(() => {
    collectionStore.collectionDetail = null
    isAIMode.value = false
  })

  const keys = useMagicKeys({
    passive: false,
    onEventFired(e) {
      if ((e.ctrlKey || e.metaKey) && e.key === 's' && e.type === 'keydown')
        e.preventDefault()
    },
  })

  whenever(keys.cmd_s, () => {
    if (readonly.value || localStorage.getItem('apicat.com.save.tip'))
      return
    ElMessage.closeAll()
    ElMessage({
      showClose: true,
      duration: 0,
      message: t('app.tips.autoSave'),
      onClose() {
        localStorage.setItem('apicat.com.save.tip', '1')
      },
    })
  })

  return {
    collection,
    isLoading,
    isSaving,
    isSaveError,
    isShowAIStyle,
    isShowAIStyleForTitle,

    handleEditorEvent,
    handleTitleBlur,
    handleContentUpdate,
    onShareCollectionBtnClick,
    onExportCollectionBtnClick,
    goCollectionhistory,

    readonly,
    toggleMode,
    switchToWriteMode,
    titleInputRef: inputRef,
    onDocumentLayoutClick,
  }
}
