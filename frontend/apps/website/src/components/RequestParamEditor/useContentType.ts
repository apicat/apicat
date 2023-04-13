import { useNodeAttrs, HTTP_REQUEST_NODE_KEY } from '@/hooks/useNodeAttrs'
import { RequestContentTypesMap } from '@/commons'
import { createDefaultSchema } from '@/views/document/components/createDefaultDefinition'

export const useContentType = (props: any) => {
  const nodeAttrs = useNodeAttrs(props, HTTP_REQUEST_NODE_KEY)

  // 缓存各个content的值
  const contentValues: any = ref({
    [RequestContentTypesMap.none]: {},
    [RequestContentTypesMap['form-data']]: { schema: createDefaultSchema() },
    [RequestContentTypesMap['x-www-form-urlencoded']]: { schema: createDefaultSchema() },
    [RequestContentTypesMap.json]: { schema: createDefaultSchema() },
    [RequestContentTypesMap.xml]: { schema: createDefaultSchema() },
    [RequestContentTypesMap.raw]: { schema: createDefaultSchema() },
    [RequestContentTypesMap.binary]: { schema: createDefaultSchema() },
  })

  const currentContentTypeRef = ref(RequestContentTypesMap.none)

  const initContentValueAndType = () => {
    if (!nodeAttrs.value.content) {
      nodeAttrs.value.content = {}
    }

    const keys = Object.keys(RequestContentTypesMap)
    for (let i = 0; i < keys.length; i++) {
      const key = (RequestContentTypesMap as any)[keys[i]]
      // 多个content，默认取一个
      if (Reflect.has(nodeAttrs.value.content, key)) {
        currentContentTypeRef.value = key
        contentValues.value[key] = nodeAttrs.value.content[key]
        break
      }
    }
  }

  const setNodeAttrsContentValue = (val?: any) => {
    nodeAttrs.value.content = {
      [currentContentTypeRef.value]: toRaw(val || contentValues.value[currentContentTypeRef.value]),
    }
  }

  const bodyCount = computed(() => (currentContentTypeRef.value !== RequestContentTypesMap.none ? 1 : 0))

  const handleChooseFile = (file: File) => {
    if (file) {
      const schema = contentValues.value[currentContentTypeRef.value].schema
      schema.example = file.name
      setNodeAttrsContentValue()
    }
  }

  // update:modelValue
  watch(currentContentTypeRef, () => setNodeAttrsContentValue())
  watch(contentValues, () => setNodeAttrsContentValue(), { deep: true })
  watch(nodeAttrs, () => initContentValueAndType(), { immediate: true })

  return {
    RequestContentTypesMap,

    currentContentTypeRef,
    contentValues,

    bodyCount,

    handleChooseFile,
  }
}
