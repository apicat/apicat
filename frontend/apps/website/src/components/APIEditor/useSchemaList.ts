import { ElMessage } from 'element-plus'
import { APICatSchemaObjectCustom } from './types'
import { debounce } from 'lodash-es'

export const useSchemaList = (
  props: any,
  emit: any,
  transformModel: (_m: APICatSchemaObjectCustom[]) => unknown,
  onParamNameValid?: (oldName: string, newName: string) => void
) => {
  const { onCreate, onDelete, onChange } = props

  const newname = ref('')

  const model: Ref<APICatSchemaObjectCustom[]> = ref([])

  const changeNotify = debounce((item?: APICatSchemaObjectCustom) => {
    // update
    if (item) {
      const { _name, ...others } = item
      onChange && onChange(toRaw(others))
    }
    transformModel && emit('update:modelValue', transformModel(model.value))
  }, 500)

  const validParamName = (v: string, item?: APICatSchemaObjectCustom) => {
    if (v == '') {
      ElMessage.error('参数名不能为空')
      return false
    }
    if (model.value.find((item) => item.name == v)) {
      ElMessage.error(`参数「${v}」重复`)
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
    onParamNameValid && onParamNameValid(item.name, v)
    item.name = v
    onChange && onChange(item)
    changeNotify()
  }, 500)

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
