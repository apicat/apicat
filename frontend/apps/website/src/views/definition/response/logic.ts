import { debounce, isEmpty } from 'lodash-es'
import { storeToRefs } from 'pinia'
import useDefinitionStore from '@/store/definition'
import { getDefinitionResponseDetail } from '@/api/definitionResponse'
import { DefinitionResponse } from '@/typings'
import { useParams } from '@/hooks/useParams'
import { ElMessage } from 'element-plus'
import { useI18n } from 'vue-i18n'
import useDefinitionResponseStore from '@/store/definitionResponse'

export const useDefinitionResponseLogic = () => {
  const route = useRoute()
  const definitionStore = useDefinitionStore()
  const { definitions } = storeToRefs(definitionStore)

  const { project_id } = useParams()

  const [isLoading, getDefinitionResponseDetailRequest] = getDefinitionResponseDetail()

  const hasDocument = ref(true)
  const response = ref<DefinitionResponse | null>(null)

  const getDetail = async () => {
    const id = parseInt(route.params.response_id as string, 10)

    if (isNaN(id)) {
      hasDocument.value = false
      return
    }
    hasDocument.value = true

    try {
      response.value = await getDefinitionResponseDetailRequest({ project_id, id })
    } catch (error) {
      //
    }
  }

  watch(
    () => route.params.response_id,
    async () => await getDetail(),
    { immediate: true }
  )

  definitionStore.$onAction(({ name, after }) => {
    if (name === 'deleteDefinition') {
      after(() => getDetail())
    }
  })

  return {
    isLoading,
    hasDocument,
    definitionSchemas: definitions,
    response,
  }
}

export const useEditDefinitionResponseLogic = () => {
  const { response, ...others } = useDefinitionResponseLogic()

  const { t } = useI18n()
  const { project_id } = useParams()
  const definitionResponseStore = useDefinitionResponseStore()

  const isSaving = ref(false)

  const definitionResponseTree: any = inject('definitionResponseTree')

  watch(
    response,
    debounce(async (newVal, oldVal) => {
      if (!oldVal || !newVal || !oldVal.id || newVal.id !== oldVal.id) {
        return
      }

      if (isEmpty(newVal.name)) {
        ElMessage.error(t('app.definitionResponse.form.title'))
        return
      }

      isSaving.value = true
      try {
        const data: any = unref(response)
        await definitionResponseStore.updateDefinition({ project_id, ...data })
        definitionResponseTree.updateTitle(data.id, newVal.name)
      } catch (e) {
        //
      } finally {
        isSaving.value = false
      }
    }, 300),
    { deep: true }
  )

  return {
    ...others,
    response,
    isSaving,
  }
}
