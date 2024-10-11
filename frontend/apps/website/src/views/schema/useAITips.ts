import { debounce } from 'lodash-es'
import type { JSONSchemaTable } from '@apicat/components'
import axios, { AxiosError } from 'axios'
import { apiGetAIModel } from '@/api/project/definition/schema'

export function useAITips(project_id: string, schema: Ref<Definition.SchemaNode | null>, readonly: Ref<boolean>, updateSchema: (projectID: string, schema: Definition.SchemaNode) => Promise<void | undefined>) {
  const { escape } = useMagicKeys()
  const jsonSchemaTableIns = ref<InstanceType<typeof JSONSchemaTable>>()

  const schemaName = ref('')
  const isAIMode = ref(false)
  const isShowAIStyle = ref(false)
  const isLoadingAICollection = ref(false)
  const preSchema = ref<Definition.SchemaNode | null>(null)

  // 避免请求后，文档不匹配问题
  const requestID = ref<string>()
  let abortController: AbortController | null = null

  // 不允许AI提示系列操作判断条件
  const notAllowAITips = () => preSchema.value?.id !== schema.value?.id || !isAIMode.value || !schema.value || readonly.value

  // 获取AI提示数据
  async function getAITips() {
    if (notAllowAITips())
      return

    // 取消上次请求
    abortController?.abort()

    requestID.value = `${Date.now()},${schema.value!.id}`

    try {
      abortController = new AbortController()
      isLoadingAICollection.value = true
      const res = await apiGetAIModel(project_id, { requestID: unref(requestID), modelID: schema.value!.id, title: schema.value!.name }, { signal: abortController.signal })
      const { schema: aiSchema, requestID: resRequestID } = res || {}
      if (requestID.value === resRequestID && aiSchema) {
        schema.value!.schema = aiSchema
        isShowAIStyle.value = true
      }

      abortController = null
      // 重置请求标识
      isLoadingAICollection.value = false
      requestID.value = ''
    }
    catch (error: any) {
      // Cancelled Error 不需要重置
      if (axios.isCancel(error)) {
        isLoadingAICollection.value = false
        requestID.value = ''
      }
    }
  }

  // 标题失去焦点时,延迟600避免title&path的debounce冲突
  function handleTitleOrPathBlur() {
    // 获取AI数据中
    if (isLoadingAICollection.value) {
      cancelAITips()
      return
    }

    confirmAITips()
  }

  // 取消AI提示
  function cancelAITips() {
    // 重置请求ID，避免请求后，文档不匹配问题
    requestID.value = ''
    isShowAIStyle.value = false
    isLoadingAICollection.value = false
    abortController?.abort()
    // 还原文档
    if (preSchema.value && schema.value)
      schema.value.schema = JSON.parse(JSON.stringify(preSchema.value.schema))
  }

  // 确认AI提示
  function confirmAITips() {
    if (notAllowAITips())
      return
    // trigger watch name
    schema.value!.name = schema.value!.name
    isShowAIStyle.value = false
    isAIMode.value = false
    try {
      const copySchemaStr = JSON.stringify(schema.value)
      const copyPreSchemaStr = JSON.stringify(preSchema.value)
      if (copySchemaStr === copyPreSchemaStr)
        return

      updateSchema(project_id, schema.value!)
      // 保存历史文档
      preSchema.value = JSON.parse(copySchemaStr)
    }
    catch (e) {
      console.error('confirmAITips error', e)
      preSchema.value = null
    }
  }

  watch(schemaName, async () => {
    if (isAIMode.value || jsonSchemaTableIns.value?.isEmpty()) {
      isAIMode.value = true
      await getAITips()
    }
  })

  whenever(escape, () => {
    if (readonly.value || !isAIMode.value || !schema.value)
      return
    cancelAITips()
  })

  return {
    isAIMode,
    isShowAIStyle,
    jsonSchemaTableIns,
    schemaName,
    preSchema,

    handleTitleBlur: handleTitleOrPathBlur,
  }
}
