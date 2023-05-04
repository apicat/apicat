import { ElMessage } from 'element-plus'
import { APICatSchemaObjectCustom } from './types'
import { debounce } from 'lodash-es'
import { useI18n } from 'vue-i18n'

export const useSchemaList = (
  props: any,
  emit: any,
  transformModel: (_m: APICatSchemaObjectCustom[]) => unknown,
  onChangeParamNameSuccess?: (oldName: string, newName: string) => void
) => {
  const { onCreate, onDelete, onChange } = props
  const { t } = useI18n()

  const newname = ref('')

  const model: Ref<APICatSchemaObjectCustom[]> = ref([])

  const changeNotify = debounce((item?: APICatSchemaObjectCustom) => {
    // update
    if (item) {
      const { _name, ...others } = item
      onChange && onChange(toRaw(others))
    }
    transformModel && emit('update:modelValue', transformModel(model.value))
  }, 300)

  const validParamName = (v: string, item?: APICatSchemaObjectCustom) => {
    if (v == '') {
      ElMessage.error(t('editor.common.error.emptyParamName'))
      return false
    }
    if (model.value.find((item) => item.name == v)) {
      ElMessage.error(t('editor.common.error.paramNameDuplicate', [v]))
      if (item) {
        item._name = item.name
      }
      return false
    }
    return true
  }

  const onParamNameChange = debounce((item: APICatSchemaObjectCustom, v: string) => {
    if (!validParamName(v, item)) {
      // item._name = ''
      return
    }
    onChangeParamNameSuccess && onChangeParamNameSuccess(item.name, v)
    item.name = v
    onChange && onChange(item)
    changeNotify()
  }, 200)

  const addHandler = async (v: string) => {
    if (!validParamName(v)) {
      return
    }

    let newItem: APICatSchemaObjectCustom = {
      name: v,
      required: false,
      schema: { type: 'string' },
    }

    try {
      newname.value = ''

      if (onCreate) {
        const data = await onCreate({ ...newItem })
        newItem = { ...newItem, ...data }
      }

      newItem._name = v
      model.value.push(newItem)
      changeNotify()
    } catch (error) {
      //
    }
  }

  const delHandler = (i: number) => {
    const deleteItem = model.value.splice(i, 1)
    if (deleteItem && deleteItem.length && onDelete) {
      onDelete && onDelete({ ...deleteItem[0] })
    }
    changeNotify()
  }

  return {
    newname,
    model,
    onParamNameChange,
    addHandler,
    delHandler,
    changeNotify,
  }
}
