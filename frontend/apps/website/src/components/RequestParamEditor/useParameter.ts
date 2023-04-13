import { useNodeAttrs, HTTP_REQUEST_NODE_KEY } from '@/hooks/useNodeAttrs'

export const useParameter = (props: any, propKey?: string) => {
  const nodeAttrs = useNodeAttrs(props, HTTP_REQUEST_NODE_KEY, propKey)

  const headers = computed({
    get: () => nodeAttrs.value.parameters.header,
    set: (val: any) => {
      nodeAttrs.value.parameters.header = val
    },
  })

  const cookies = computed({
    get: () => nodeAttrs.value.parameters.cookie,
    set: (val: any) => {
      nodeAttrs.value.parameters.cookie = val
    },
  })

  const queries = computed({
    get: () => nodeAttrs.value.parameters.query,
    set: (val: any) => {
      nodeAttrs.value.parameters.query = val
    },
  })

  const paths = computed({
    get: () => nodeAttrs.value.parameters.path,
    set: (val: any) => {
      nodeAttrs.value.parameters.path = val
    },
  })

  const headersCount = computed(() => headers.value.length || '')
  const cookiesCount = computed(() => cookies.value.length || '')
  const queriesCount = computed(() => queries.value.length || '')
  const pathsCount = computed(() => paths.value.length || '')

  return {
    headers,
    cookies,
    queries,
    paths,

    headersCount,
    cookiesCount,
    queriesCount,
    pathsCount,
  }
}
