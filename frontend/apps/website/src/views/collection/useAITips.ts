import { debounce } from 'lodash-es'
import { EDITOR_NODE_EVENT } from '@apicat/editor'
import { AxiosError } from 'axios'
import { apiGetAICollection } from '@/api/project/collection'

export function useAITips(project_id: string, collection: Ref<CollectionAPI.ResponseCollectionDetail | null>, readonly: Ref<boolean>, updateCollection: (projectID: string, collection: CollectionAPI.ResponseCollectionDetail) => Promise<void>) {
  const { escape } = useMagicKeys()

  const isAIMode = ref(false)
  const isShowAIStyle = ref(false)
  const isShowAIStyleForTitle = ref(false)
  const isLoadingAICollection = ref(false)
  const preCollection = ref<CollectionAPI.ResponseCollectionDetail | null>(null)

  const docTitle = ref('')
  const docPath = ref('')

  // 避免请求后，文档不匹配问题
  const requestID = ref<string>()
  let abortController: AbortController | null = null

  // 不允许AI提示系列操作判断条件
  const notAllowAITips = () => !isAIMode.value || !collection.value || readonly.value

  // 获取AI提示数据
  async function getAITips(callback?: () => void) {
    const path = collection.value?.content?.[0]?.attrs?.path
    const title = collection.value?.title

    if (!title || notAllowAITips())
      return
    // 取消上次请求
    abortController?.abort()

    requestID.value = `${Date.now()},${collection.value!.id}`
    try {
      abortController = new AbortController()
      isLoadingAICollection.value = true
      const { content: aiCollection, requestID: resRequestID } = await apiGetAICollection(project_id, { requestID: unref(requestID), title, path }, { signal: abortController.signal })
      if (requestID.value === resRequestID && isLoadingAICollection.value) {
        collection.value!.content = aiCollection.content
        collection.value!.title = aiCollection.title
        callback && callback()
      }

      abortController = null
      // 重置请求标识
      isLoadingAICollection.value = false
      requestID.value = ''
    }
    catch (error: any) {
      console.error(error)
      // Cancelled Error 不需要重置
      if (error && AxiosError.ERR_CANCELED !== error.code) {
        isLoadingAICollection.value = false
        requestID.value = ''
      }
    }
  }

  // 编辑器事件处理
  function handleEditorEvent(event: string, data: any = {}) {
    // 当path改变时，获取AI数据
    if (event === EDITOR_NODE_EVENT.HTTP_PATH_CHANGE) {
      docPath.value = data.path
      return
    }

    // 当path失焦时
    if (event === EDITOR_NODE_EVENT.HTTP_PATH_BLUR)
      handleTitleOrPathBlur()
  }

  // 标题失去焦点时,延迟600避免title&path的debounce冲突
  const handleTitleOrPathBlur = debounce(() => {
    // 获取AI数据中
    if (isLoadingAICollection.value) {
      cancelAITips()
      return
    }

    confirmAITips()
  }, 600)

  // 取消AI提示
  function cancelAITips() {
    // 重置请求ID，避免请求后，文档不匹配问题
    requestID.value = ''
    isShowAIStyle.value = false
    isShowAIStyleForTitle.value = false
    isLoadingAICollection.value = false
    abortController?.abort()
    // 还原文档
    if (preCollection.value && collection.value)
      collection.value.content = JSON.parse(JSON.stringify(preCollection.value)).content
  }

  // 确认AI提示
  function confirmAITips() {
    if (notAllowAITips())
      return
    // trigger watch title
    collection.value!.title = collection.value!.title
    isShowAIStyle.value = false
    isShowAIStyleForTitle.value = false
    try {
      const copyCollectionStr = JSON.stringify(collection.value)
      const copyPreCollectionStr = JSON.stringify(preCollection.value)
      if (copyCollectionStr === copyPreCollectionStr)
        return

      updateCollection(project_id, collection.value!)
      // 保存历史文档
      preCollection.value = JSON.parse(copyCollectionStr)
    }
    catch (e) {
      console.error('confirmAITips error', e)
      preCollection.value = null
    }
  }

  watch(docTitle, async () => {
    await getAITips(() => {
      isShowAIStyleForTitle.value = true
      isShowAIStyle.value = false
    })
  })

  watch(docPath, async () => {
    await getAITips(() => {
      isShowAIStyle.value = true
      isShowAIStyleForTitle.value = false
    })
  })

  whenever(escape, () => {
    if (readonly.value || !isAIMode.value || !collection.value)
      return
    cancelAITips()
  })

  return {
    isAIMode,
    isShowAIStyle,
    isShowAIStyleForTitle,
    isLoadingAICollection,

    docTitle,
    docPath,
    preCollection,

    handleEditorEvent,
    handleTitleBlur: handleTitleOrPathBlur,
  }
}
