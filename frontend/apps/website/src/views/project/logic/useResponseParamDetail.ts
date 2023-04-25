import { getResponseParam } from '@/api/commonResponse'
import useCommonResponseStore from '@/store/commonResponse'
import { ProjectInfo, APICatCommonResponseCustom } from '@/typings'
import { ElMessage } from 'element-plus'

export const useResponseParamDetail = ({ id: project_id }: Pick<ProjectInfo, 'id'>) => {
  const commonResponseStore = useCommonResponseStore()

  const handleExpand = async (isExpand: boolean, item: APICatCommonResponseCustom) => {
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

  const handleSubmit = async (param: APICatCommonResponseCustom) => {
    const { detail } = param
    if (!detail) {
      return
    }

    if (!detail.name) {
      ElMessage.error('响应名称不能为空')
      return
    }

    param.isLoading = true
    try {
      // 更新
      if (detail.id) {
        await commonResponseStore.updateResponseParam(project_id, detail)
      }

      // 添加
      if (!detail.id) {
        const responseParamDetail = await commonResponseStore.addCommonResponse(project_id, detail)
        param.isLocal = false
        param.id = responseParamDetail.id
        param.detail = responseParamDetail
      }
    } catch (e) {
      //
    } finally {
      param.isLoading = false
    }
  }

  return {
    handleExpand,
    handleSubmit,
  }
}
