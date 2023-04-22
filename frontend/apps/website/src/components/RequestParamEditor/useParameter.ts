import { useNodeAttrs, HTTP_REQUEST_NODE_KEY } from '@/hooks/useNodeAttrs'
import uesGlobalParametersStore from '@/store/globalParameters'
import { storeToRefs } from 'pinia'

export const useParameter = (props: any, propKey?: string) => {
  const nodeAttrs = useNodeAttrs(props, HTTP_REQUEST_NODE_KEY, propKey)
  const globalParametersStore = uesGlobalParametersStore()
  const { parameters: globalParameters } = storeToRefs(globalParametersStore)

  globalParametersStore.$subscribe((e, state) => {
    console.log(e, state)
  })

  const headers = computed({
    get: () => nodeAttrs.value.parameters.header,
    set: (val: any) => {
      nodeAttrs.value.parameters.header = val
    },
  })

  const globalHeaders = computed(() =>
    globalParameters.value.header.map((param) => {
      return { ...param, required: param.required ? '是' : '否', isUse: !(nodeAttrs.value.globalExcepts.header || []).includes(param.id) }
    })
  )

  const switchGlobalParameter = (id: string | number, isUse: boolean, _in: string) => {
    nodeAttrs.value.globalExcepts[_in] = isUse ? nodeAttrs.value.globalExcepts[_in].filter((item: string | number) => item !== id) : [...nodeAttrs.value.globalExcepts[_in], id]
  }

  const switchGlobalHeader = (id: any, isUse: any) => {
    console.log(id, isUse)
    switchGlobalParameter(id, isUse, 'header')
  }
  const switchGlobalCookie = (id: string | number, isUse: boolean) => switchGlobalParameter(id, isUse, 'cookie')
  const switchGlobalQuery = (id: string | number, isUse: boolean) => switchGlobalParameter(id, isUse, 'query')
  const switchGlobalPath = (id: string | number, isUse: boolean) => switchGlobalParameter(id, isUse, 'path')

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

  const headersCount = computed(() => headers.value.length + globalHeaders.value.length || '')
  const cookiesCount = computed(() => cookies.value.length || '')
  const queriesCount = computed(() => queries.value.length || '')
  const pathsCount = computed(() => paths.value.length || '')

  return {
    headers,
    globalHeaders,
    cookies,
    queries,
    paths,

    headersCount,
    cookiesCount,
    queriesCount,
    pathsCount,

    switchGlobalHeader,
  }
}
