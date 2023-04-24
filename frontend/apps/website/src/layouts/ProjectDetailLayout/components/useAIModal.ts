export const useAIModal = (initTree: any) => {
  const aiPromptModalRef = shallowRef()

  const onCreateSuccess = async (active_id: any) => {
    initTree(active_id)
  }
  return {
    aiPromptModalRef,
    onCreateSuccess,
  }
}
