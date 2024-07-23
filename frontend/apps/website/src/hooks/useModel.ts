import type { FormInstance } from 'element-plus'

export function useModal(formRef?: Ref<FormInstance>): { dialogVisible: Ref<boolean>, showModel: () => void, hideModel: () => void } {
  const dialogVisible = ref(false)

  watch(dialogVisible, () => {
    if (!dialogVisible.value) {
      const formIns = unref(formRef)
      if (formIns) {
        formIns.resetFields()
        formIns.clearValidate()
      }
    }
  })

  const showModel = () => {
    dialogVisible.value = true
  }

  const hideModel = () => {
    dialogVisible.value = false
  }

  return {
    dialogVisible,
    showModel,
    hideModel,
  }
}
