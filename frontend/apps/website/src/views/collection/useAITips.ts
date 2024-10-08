import { EDITOR_NODE_EVENT } from '@apicat/editor'
import { apiGetAICollection } from '@/api/project/collection'

export function useAITips(project_id: string, collection: Ref<CollectionAPI.ResponseCollectionDetail | null>, readonly: Ref<boolean>) {
  const isAIMode = ref(false)
  const isShowAIStyle = ref(false)
  const isShowAIStyleForTitle = ref(false)
  const isLoadingAICollection = ref(false)
  const preCollection = ref<CollectionAPI.ResponseCollectionDetail | null>(null)

  const docTitle = ref('')
  const docPath = ref('')

  // 避免请求后，文档不匹配问题
  const requestID = ref<string>()

  async function getAITips(callback?: () => void) {
    const path = collection.value?.content?.[0]?.attrs?.path
    if (!collection.value?.title || !path || !isAIMode.value || !collection.value || isLoadingAICollection.value || readonly.value)
      return

    requestID.value = `${Date.now()},${collection.value.id}`

    try {
      isLoadingAICollection.value = true
      const { content: aiCollection, requestID: resRequestID } = await apiGetAICollection(project_id, { requestID: unref(requestID), title: unref(docTitle), path: unref(docPath) })
      isLoadingAICollection.value = false
      if (requestID.value === resRequestID) {
        collection.value.content = aiCollection.content
        collection.value.title = aiCollection.title
        callback && callback()
      }
    }
    catch (error) {
      console.error(error)
      isLoadingAICollection.value = false
    }
  }

  function handleEditorEvent(event: string, data: any = {}) {
    // 当path改变时，获取AI数据
    if (event === EDITOR_NODE_EVENT.HTTP_PATH_CHANGE) {
      docPath.value = data.path
      return
    }

    // 当path失焦时
    if (event === EDITOR_NODE_EVENT.HTTP_PATH_BLUR) {
      //
    }
  }

  // 标题失去焦点时
  function handleTitleBlur() {
    //
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

  return {
    isAIMode,
    isShowAIStyle,
    isShowAIStyleForTitle,
    isLoadingAICollection,

    docTitle,
    docPath,
    preCollection,

    handleEditorEvent,
    handleTitleBlur,
  }
}
