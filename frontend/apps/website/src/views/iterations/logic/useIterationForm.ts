import { storeToRefs } from 'pinia'
import type { FormInstance } from 'element-plus'
import { useI18n } from 'vue-i18n'
import type { IterationFormEmits, IterationFormProps } from '../components/IterationForm.vue'
import { useTeamStore } from '@/store/team'
import { apiCreateIteration, apiEditIterationInfo, apiGetIterationInfo } from '@/api/iteration/index'
import useApi from '@/hooks/useApi'

export type CreateOrEditIteration = GlobalAPI.Merge<
  IterationAPI.RequestCreateIteration,
  IterationAPI.RequestEditIteration
>

export function useIterationForm(props: IterationFormProps, emits: IterationFormEmits) {
  const { t } = useI18n()
  const defaultCreateIteration: CreateOrEditIteration = {
    title: '',
    projectID: '',
    description: '',
    collectionIDs: [],
  }
  const [isLoadedIteration, getIterationInfo] = useApi(apiGetIterationInfo)
  const iterationFormRef = shallowRef()
  const isLoadingForSubmit = ref(false)
  const iterationInfo = ref<CreateOrEditIteration>({
    ...defaultCreateIteration,
  })
  const iterationRules = {
    title: [{ required: true, message: t('app.iteration.form.inpIterNameTip') }],
    projectID: [{ required: true, message: t('app.iteration.form.selectProjectTip') }],
    description: [{ message: t('app.iteration.form.descTip') }],
  }

  const { iterationID: iterationIDRef } = toRefs(props)
  const { currentID } = storeToRefs(useTeamStore())
  const isEditMode = computed(() => iterationIDRef.value !== null)

  async function handleSubmit(formIns: FormInstance) {
    try {
      await formIns.validate()
      isLoadingForSubmit.value = true
      if (isEditMode.value) {
        const { projectID, ...data } = iterationInfo.value
        await apiEditIterationInfo(iterationIDRef.value!, data)
      }
      else {
        await apiCreateIteration(currentID.value, iterationInfo.value)
      }
      handleResetForm()
      emits('success')
    }
    catch (error) {
      //
    }
    finally {
      isLoadingForSubmit.value = false
    }
  }

  function handleResetForm() {
    iterationInfo.value = { ...defaultCreateIteration }
    iterationFormRef.value?.resetFields()
  }

  function handleCancel() {
    emits('cancel')
  }

  watch(
    iterationIDRef,
    async () => {
      // 没有迭代ID
      if (!iterationIDRef.value) {
        handleResetForm()
        return
      }

      try {
        const iteration = (await getIterationInfo(iterationIDRef.value))!
        iterationInfo.value = {
          title: iteration.title || '',
          projectID: iteration.project?.id || '',
          description: iteration.description,
          collectionIDs: [],
        }
      }
      catch (error) {
        handleResetForm()
      }
    },
    { immediate: true },
  )

  return {
    isEditMode,
    isLoadedIteration,
    isLoadingForSubmit,

    iterationFormRef,
    iterationIDRef,
    iterationInfo,
    iterationRules,

    handleSubmit,
    handleCancel,
  }
}
