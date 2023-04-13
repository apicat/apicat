import { HttpDocument } from '@/typings'

export const HTTP_REQUEST_NODE_KEY = 'apicat-http-request'
export const HTTP_RESPONSE_NODE_KEY = 'apicat-http-response'
export const HTTP_URL_NODE_KEY = 'apicat-http-url'

export const useNodeAttrs = (props: any, filterKey: string, propsKey?: string, defaultValue?: any) => {
  const doc: Ref<HttpDocument> = useVModel(props, propsKey ?? 'modelValue', undefined, { passive: true, defaultValue })

  const nodeAttrs = computed({
    get: () => {
      const node = doc.value.content.find((item: any) => item.type === filterKey)
      return node.attrs
    },
    set: (val) => {
      const node = doc.value.content.find((item: any) => item.type === filterKey)
      node.attrs = val
    },
  })

  return nodeAttrs
}
