import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { storeToRefs } from 'pinia'
import { delay } from '@apicat/shared'
import type { PageModeCtx } from '../composables/usePageMode'
import useApi from '@/hooks/useApi'
import { useCollectionsStore } from '@/store/collections'
import { useProjectLayoutContext } from '@/layouts/ProjectDetailLayout'
import { getCollectionHistoryPath } from '@/router/history'
import { injectPagesMode } from '@/layouts/ProjectDetailLayout/composables/usePagesMode'
import { useTitleInputFocus } from '@/hooks/useTitleInputFocus'
import { injectAsyncInitTask } from '@/hooks/useWaitAsyncTask'
import { apiGetAICollection, isEmptyContent } from '@/api/project/collection'

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

  // 是否是AI模式，避免AI推理数据被保存
  const isAIMode = ref(false)

  function handleTitleBlur() {
    const title = collection.value?.title || ''
    if (!title || !title.trim()) {
      collection.value!.title = oldTitle
      oldTitle = ''
    }
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

      if (!isAIMode.value)
        await updateCollection(props.project_id, n)
      // else
        // await apiGetAICollection(props.project_id, { requestID: `${Date.now()}`, title: n.title, path: '/api/login' })
    }
  }, { debounce: 200 })

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
        await task
        oldTitle = collection.value?.title || ''
        if (!readonly.value)
          focus()
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
    isAIMode,

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
