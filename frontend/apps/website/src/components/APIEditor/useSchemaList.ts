import { ElMessage } from 'element-plus'
import { APICatSchemaObject } from './types'
import { debounce } from 'lodash-es'

export type APICatSchemaObjectCustom = APICatSchemaObject & { _name?: string }

export const useSchemaList = (emit: any, transformModel: (_m: APICatSchemaObjectCustom[]) => unknown, onParamNameValid?: (oldName: string, newName: string) => void) => {
  const newname = ref('')

  const model: Ref<APICatSchemaObjectCustom[]> = ref([])

  const changeNotify = () => {
    transformModel && emit('update:modelValue', transformModel(model.value))
  }

  const validParamName = (v: string) => {
    if (v == '') {
      ElMessage.error('参数名不能为空')
      return false
    }
    if (model.value.find((item) => item.name == v)) {
      ElMessage.error(`参数「${v}」重复`)
      return false
    }
    return true
  }

  const onParamNameChange = debounce((item: APICatSchemaObjectCustom, v: string) => {
    if (!validParamName(v)) {
      // item._name = ''
      return
    }

    onParamNameValid && onParamNameValid(item.name, v)
    item.name = v
    changeNotify()
  }, 200)

  const addHandler = (v: string) => {
    if (!validParamName(v)) {
      return
    }

    newname.value = ''
    model.value.push({
      name: v,
      _name: v,
      schema: { type: 'string' },
    })

    changeNotify()
  }

  const delHandler = (i: number) => {
    model.value.splice(i, 1)
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
