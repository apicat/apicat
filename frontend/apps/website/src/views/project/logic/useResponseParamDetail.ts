// import { getResponseParam } from '@/api/commonResponse'
import useCommonResponseStore from '@/store/commonResponse'
import { ProjectInfo, APICatCommonResponseCustom } from '@/typings'
import { ElMessage } from 'element-plus'
import { useI18n } from 'vue-i18n'

export const useResponseParamDetail = ({ id: project_id }: Pick<ProjectInfo, 'id'>) => {
  const commonResponseStore = useCommonResponseStore()
  const { t } = useI18n()

  const handleExpand = async (isExpand: boolean, item: APICatCommonResponseCustom) => {
    // if (isExpand && !item.isLoaded) {
    //   item.isLoading = true
    //   try {
    //     const responseParamDetail = await getResponseParam({ project_id, response_id: item.id })
    //     item.isLoaded = true
    //     item.detail = responseParamDetail
    //   } finally {
    //     item.isLoading = false
    //   }
    // }
  }

  const handleSubmit = async (param: APICatCommonResponseCustom) => {
    const { detail } = param
    if (!detail) {
      return
    }

    if (!detail.name) {
      ElMessage.error(t('app.response.rules.name'))
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
        const responseParamDetail: any = await commonResponseStore.addCommonResponse(project_id, detail)
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
