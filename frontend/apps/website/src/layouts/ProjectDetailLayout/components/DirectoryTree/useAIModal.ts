import AcTree from '@/components/AcTree'

export const useAIModal = (initDocumentTree: any) => {
  const aiPromptModalRef = shallowRef()

  const onCreateSuccess = async (doc_id: any) => {
    initDocumentTree(doc_id)
  }
  return {
    aiPromptModalRef,
    onCreateSuccess,
  }
}
