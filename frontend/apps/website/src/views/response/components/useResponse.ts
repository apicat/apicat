/* eslint-disable no-unreachable-loop */
import debounce from 'lodash-es/debounce'
import type { JSONSchema, ResponseContentType } from '@apicat/editor'
import { watchDebounced } from '@vueuse/core'
import type { ResponseFormEmits, ResponseFormProps } from './ResponseForm.vue'

export function useResponse(props: ResponseFormProps, emits?: ResponseFormEmits) {
  const definitionResponseRef = toRef(props, 'response')

  function update() {
    emits?.('update:response', unref(definitionResponseRef))
  }

  const debounceUpdate = debounce(update, 200)

  const {
    contentType: resContentType,
    schema: resSchema,
    examples: resExamples,
  } = parseResponse(unref(definitionResponseRef))

  const contentSchema = ref<JSONSchema>(resSchema)
  const contentType = ref<string>(resContentType)
  const examples = ref<Record<string, any>>(resExamples)
  const headers = computed({
    get: () => unref(definitionResponseRef).header || [],
    set: (_headers: any) => {
      const response = unref(definitionResponseRef)
      response.header = _headers
      debounceUpdate()
    },
  })

  const isJSONSchema = computed(() => contentType.value === 'application/json' || contentType.value === 'application/xml')

  watchDebounced(
    [examples, contentType, contentSchema],
    ([, newContentType], [, oldContentType]) => {
      const response = unref(definitionResponseRef)
      if (response.content) {
        const content = response.content as any
        delete content[oldContentType]
        content[newContentType] = {
          schema: isJSONSchema.value ? contentSchema.value : {},
          examples: examples.value,
        }
      }
      update()
    },
    {
      debounce: 200,
    },
  )

  watch(definitionResponseRef, () => {
    const {
      contentType: resContentType,
      schema: resSchema,
      examples: resExamples,
    } = parseResponse(unref(definitionResponseRef))
    contentType.value = resContentType
    contentSchema.value = resSchema
    examples.value = resExamples
  })

  // 解析response
  function parseResponse(response?: Definition.ResponseDetail): {
    schema: JSONSchema
    contentType: ResponseContentType
    examples: Record<string, any>
  } {
    const contentType = getContentTypeByResponse(response)
    let schema: JSONSchema = {}
    let examples: Record<string, any> = {}

    if (response && response.content) {
      const content = response.content[contentType]! || {}
      schema = content.schema || {}
      examples = content.examples || {}
    }

    return {
      contentType,
      schema,
      examples,
    }
  }

  function getContentTypeByResponse(response?: Definition.ResponseDetail): ResponseContentType {
    for (const x in response?.content) return x as ResponseContentType
    return 'application/json'
  }

  return {
    contentType,
    contentSchema,
    isJSONSchema,
    examples,
    headers,
  }
}
