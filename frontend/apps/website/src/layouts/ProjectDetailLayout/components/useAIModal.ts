export const useAIModal = (_onCreateSuccess?: any) => {
  const aiPromptModalRef = shallowRef()

  const onCreateSuccess = (active_id: any) => _onCreateSuccess(active_id)

  return {
    aiPromptModalRef,
    onCreateSuccess,
  }
}
