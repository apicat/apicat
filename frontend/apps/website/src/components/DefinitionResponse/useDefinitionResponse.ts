import { isJSONSchemaContentType } from '@/commons'

export const contentTypes: Record<string, string> = {
  'application/json': 'json',
  'application/xml': 'xml',
  'text/html': 'html',
  'text/plain': 'raw',
  'application/octet-stream': 'raw',
}

export const useDefinitionResponse = (props: any) => {
  const responseRef: any = useVModel(props, 'response', undefined, { passive: true })

  const contentDefaultType = computed(() => {
    for (let x in props.response.content) {
      return x
    }
    return 'application/json'
  })

  const isJsonSchema = computed(() => isJSONSchemaContentType(contentDefaultType.value))

  const changeContentType = (v: string) => {
    const oldtype = contentDefaultType.value
    responseRef.value.content[v] = responseRef.value.content[oldtype]
    delete responseRef.value.content[oldtype]
  }

  const examples = computed({
    get:()=>{
      return responseRef.value.content[contentDefaultType.value].examples || {}
    },
    set:(value: Record<string, any>)=>{
      responseRef.value.content[contentDefaultType.value].examples = value
    }
  })

  return {
    responseRef,
    contentDefaultType,
    isJsonSchema,
    examples,

    changeContentType,
  }
}
