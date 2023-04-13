import { addResponseParam, getResponseParam, updateResponseParam } from '@/api/param'
import { ProjectInfo, ResponseListCustom } from '@/typings'

export const useResponseParamDetail = ({ id: project_id }: Pick<ProjectInfo, 'id'>) => {
  const handleExpand = async (isExpand: boolean, item: ResponseListCustom) => {
    if (isExpand && !item.isLoaded) {
      item.isLoading = true
      try {
        const responseParamDetail = await getResponseParam({ project_id, response_id: item.id })
        item.isLoaded = true
        item.detail = responseParamDetail
      } finally {
        item.isLoading = false
      }
    }
  }

  const handleSubmit = async (param: ResponseListCustom) => {
    param.isLoading = true
    const { detail } = param
    if (!detail) {
      return
    }
    try {
      // 更新
      if (detail.id) {
        const { id: response_id, ...rest } = detail
        await updateResponseParam({ project_id, response_id, ...rest })
      }

      // 添加
      if (!detail.id) {
        const responseParamDetail = await addResponseParam({ project_id, ...detail })
        param.detail = responseParamDetail
      }
    } finally {
      param.isLoading = false
    }
  }

  return {
    handleExpand,
    handleSubmit,
  }
}
