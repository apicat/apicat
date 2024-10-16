import { debounce } from 'lodash-es'
import { EDITOR_NODE_EVENT } from '@apicat/editor'
import axios from 'axios'
import { apiGetAICollection, isEmptyContent } from '@/api/project/collection'

export function useAITips(
  project_id: string,
  collection: Ref<CollectionAPI.ResponseCollectionDetail | null>,
  readonly: Ref<boolean>,
  updateCollection: (projectID: string, collection: CollectionAPI.ResponseCollectionDetail) => Promise<void>,
  currentCollectionIDRef: Ref<string>,
) {
  const { escape, tab } = useMagicKeys()

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
        const oldPath = collection.value?.content?.find(i => i.type === 'apicat-http-url').attrs?.path
        if (!oldPath || oldPath === '/')
          collection.value!.content = aiCollection.content

        else
          collection.value!.content = [...(collection.value?.content?.filter(i => i.type === 'apicat-http-url') || []), ...(aiCollection.content?.filter((i: any) => i.type !== 'apicat-http-url') || [])]

        callback && callback()
      }

      abortController = null
      // 重置请求标识
      isLoadingAICollection.value = false
      requestID.value = ''
    }
    catch (error: any) {
      // Cancelled Error 不需要重置
      if (!axios.isCancel(error)) {
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
  const handleTitleOrPathBlur = debounce(() => cancelAITips(), 600)

  // 取消AI提示
  function cancelAITips() {
    // 取消上次请求
    abortController?.abort()

    // 切换了新的内容，不需要还原
    if (Number(currentCollectionIDRef.value) !== preCollection.value?.id)
      isAIMode.value = false

    // 还原文档
    if (isAIMode.value && preCollection.value && collection.value && Number(currentCollectionIDRef.value) !== preCollection.value?.id) {
      // 仅还原content中除了path之外的数据
      const content = JSON.parse(JSON.stringify(preCollection.value)).content
      collection.value.content = [...(collection.value.content?.filter(i => i.type === 'apicat-http-url') || []), ...(content.filter((i: any) => i.type !== 'apicat-http-url') || [])]
      isAIMode.value = false
    }

    // 重置请求ID，避免请求后，文档不匹配问题
    requestID.value = ''
    isShowAIStyle.value = false
    isShowAIStyleForTitle.value = false
    isLoadingAICollection.value = false
    preCollection.value = null
  }

  // 确认AI提示
  function confirmAITips() {
    if (notAllowAITips())
      return
    // trigger watch title
    collection.value!.title = collection.value!.title
    isShowAIStyle.value = false
    isShowAIStyleForTitle.value = false
    isAIMode.value = false

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
    if (isAIMode.value || isEmptyContent(collection.value?.content)) {
      isAIMode.value = true
      // save change
      await updateCollection(project_id, collection.value!)
      preCollection.value = JSON.parse(JSON.stringify(collection.value))

      await getAITips(() => {
        isShowAIStyleForTitle.value = true
        isShowAIStyle.value = false
      })
    }
    else {
      isAIMode.value = false
    }
  })

  watch(docPath, async () => {
    if (isAIMode.value || isEmptyContent(collection.value?.content)) {
      isAIMode.value = true
      // save change
      await updateCollection(project_id, collection.value!)
      preCollection.value = JSON.parse(JSON.stringify(collection.value))

      await getAITips(() => {
        isShowAIStyle.value = true
        isShowAIStyleForTitle.value = false
      })
    }
    else {
      isAIMode.value = false
    }
  })

  whenever(escape, () => {
    if (readonly.value || !isAIMode.value || !collection.value)
      return
    cancelAITips()
  })

  whenever(tab, () => {
    // tab触发时，获取AI数据，直接取消请求
    if (isLoadingAICollection.value)
      handleTitleOrPathBlur()

    if (readonly.value || !isAIMode.value || !collection.value)
      return

    confirmAITips()
  })

  // 点击文档区域,非编辑器内点击 -> 取消AI提示
  function onDocumentLayoutClick(e: MouseEvent) {
    // 允许点击的区域的dom path 路径含有.ac-schema-editor样式，有效点击
    if (isAIMode.value && !isLoadingAICollection.value && e.composedPath().find((el: any) => el.className?.includes('ac-editor')))
      confirmAITips()
    else
      handleTitleOrPathBlur()
  }

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
    cancelAITips,
    onDocumentLayoutClick,
  }
}
