import { deleteResponseParam, getResponseParamList } from '@/api/param'
import { ProjectInfo, ResponseListCustom, ResponseList } from '@/typings'
import { createHttpResponse } from '@/views/document/components/createHttpDocument'

export const useResponseparamList = ({ id }: Pick<ProjectInfo, 'id'>) => {
  const responseParamList: Ref<ResponseListCustom[]> = ref([])

  const [isLoading, getResponseParamListApi] = getResponseParamList()

  const extendResponseParamModel = (param?: Partial<ResponseListCustom>): ResponseListCustom => {
    return {
      _id: param?.id ?? Date.now(),
      expand: false,
      isLoaded: false,
      isLoading: false,
      ...param,
    }
  }

  const createResponseParamModel = () => {
    const response = createHttpResponse({ description: '公共响应' })
    const extemdModel = extendResponseParamModel({ code: response.code, description: response.description, isLoaded: true, expand: true })
    extemdModel.detail = response
    return extemdModel
  }

  const handleAddParam = () => {
    responseParamList.value.unshift(createResponseParamModel())
  }

  const handleDeleteParam = async (item: ResponseListCustom, index: number) => {
    const { detail } = item

    let deleteId = item.id
    if (detail && detail.id) {
      deleteId = detail.id
    }

    if (deleteId) {
      item.isLoading = true
      try {
        await deleteResponseParam({ project_id: id, response_id: deleteId })
      } finally {
        item.isLoading = false
      }
    }

    responseParamList.value.splice(index, 1)
  }

  onMounted(async () => {
    const data: ResponseList[] = await getResponseParamListApi({ project_id: id })
    const list: ResponseListCustom[] = data.map((item) => extendResponseParamModel(item))
    responseParamList.value = list
  })

  return {
    isLoading,
    getResponseParamListApi,
    responseParamList,

    handleAddParam,
    handleDeleteParam,
  }
}
