import { useNodeAttrs, HTTP_REQUEST_NODE_KEY } from '@/hooks/useNodeAttrs'
import uesGlobalParametersStore from '@/store/globalParameters'
import { storeToRefs } from 'pinia'
import { useI18n } from 'vue-i18n'

export const useParameter = (props: any, propKey?: string) => {
  const { t } = useI18n()

  const nodeAttrs = useNodeAttrs(props, HTTP_REQUEST_NODE_KEY, propKey)
  const globalParametersStore = uesGlobalParametersStore()
  const { parameters: globalParameters } = storeToRefs(globalParametersStore)

  const headers = computed({
    get: () => nodeAttrs.value.parameters.header,
    set: (val: any) => {
      nodeAttrs.value.parameters.header = val
    },
  })

  const globalHeaders = computed(() =>
    globalParameters.value.header.map((param) => {
      return { ...param, required: param.required ? t('editor.table.yes') : t('editor.table.no'), isUse: !(nodeAttrs.value.globalExcepts.header || []).includes(param.id) }
    })
  )

  const globalCookies = computed(() =>
    globalParameters.value.cookie.map((param) => {
      return { ...param, required: param.required ? t('editor.table.yes') : t('editor.table.no'), isUse: !(nodeAttrs.value.globalExcepts.cookie || []).includes(param.id) }
    })
  )

  const globalQueries = computed(() =>
    globalParameters.value.query.map((param) => {
      return { ...param, required: param.required ? t('editor.table.yes') : t('editor.table.no'), isUse: !(nodeAttrs.value.globalExcepts.query || []).includes(param.id) }
    })
  )

  const globalPaths = computed(() =>
    globalParameters.value.path.map((param) => {
      return { ...param, required: param.required ? t('editor.table.yes') : t('editor.table.no'), isUse: !(nodeAttrs.value.globalExcepts.path || []).includes(param.id) }
    })
  )

  const switchGlobalParameter = (id: string | number, isUse: boolean, _in: string) => {
    nodeAttrs.value.globalExcepts[_in] = isUse ? nodeAttrs.value.globalExcepts[_in].filter((item: string | number) => item !== id) : [...nodeAttrs.value.globalExcepts[_in], id]
  }

  const switchGlobalHeader = (id: any, isUse: any) => switchGlobalParameter(id, isUse, 'header')
  const switchGlobalCookie = (id: any, isUse: any) => switchGlobalParameter(id, isUse, 'cookie')
  const switchGlobalQuery = (id: any, isUse: any) => switchGlobalParameter(id, isUse, 'query')
  const switchGlobalPath = (id: any, isUse: any) => switchGlobalParameter(id, isUse, 'path')

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

  const headersCount = computed(() => headers.value.length + globalHeaders.value.filter((i) => i.isUse).length || '')
  const cookiesCount = computed(() => cookies.value.length + globalCookies.value.filter((i) => i.isUse).length || '')
  const queriesCount = computed(() => queries.value.length + globalQueries.value.filter((i) => i.isUse).length || '')
  const pathsCount = computed(() => paths.value.length + globalPaths.value.filter((i) => i.isUse).length || '')

  return {
    headers,
    globalHeaders,
    cookies,
    globalCookies,
    queries,
    globalQueries,
    paths,
    globalPaths,

    headersCount,
    cookiesCount,
    queriesCount,
    pathsCount,

    switchGlobalHeader,
    switchGlobalCookie,
    switchGlobalQuery,
    switchGlobalPath,
  }
}
