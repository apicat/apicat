import type ProjectShareModal from '@/views/project/components/ProjectShareModal.vue'
import type DocumentShareModal from '@/views/collection/components/DocumentShareModal.vue'
import type ExportDocumentModal from '@/views/collection/components/ExportDocumentModal.vue'
import type { Options } from '@/layouts/ProjectDetailLayout/components/AICreateDialog.vue'
import type AICreateDialog from '@/layouts/ProjectDetailLayout/components/AICreateDialog.vue'

interface ProjectDetailContext {
  handleExportDocument?: (project_id?: string, doc_id?: string | number) => void
  handleShareDocument?: (project_id: string, doc_id: string | number) => void
  handleShareProject?: (project_id: string) => void
  handleAICreateCollection: (onOK: Options['onOk'], onCancel?: Options['onCancel']) => void
  handleAICreateSchema: (onOK: Options['onOk'], onCancel?: Options['onCancel']) => void
  aiPromptDialogRef: Ref<InstanceType<typeof AICreateDialog> | undefined>
}

const ProjectDetailContextKey: InjectionKey<ProjectDetailContext> = Symbol('ProjectDetailContextKey')

export function useProjectLayoutProvider() {
  const shareProjectModalRef = ref<InstanceType<typeof ProjectShareModal>>()
  const shareDocModalRef = ref<InstanceType<typeof DocumentShareModal>>()
  const exportDocModalRef = ref<InstanceType<typeof ExportDocumentModal>>()
  const AIDialogRef = ref<InstanceType<typeof AICreateDialog>>()

  provide(ProjectDetailContextKey, {
    handleShareProject: projectID => shareProjectModalRef.value?.show(projectID),
    handleShareDocument: (projectID: string, collectionID: string | number) => shareDocModalRef.value?.show(projectID, collectionID as number),
    handleExportDocument: (projectID?: string, collectionID?: string | number) => exportDocModalRef.value?.show(projectID, collectionID),
    handleAICreateCollection: (onOk, onCancel) => AIDialogRef.value?.show({ type: 'collection', onOk, onCancel }),
    handleAICreateSchema: (onOk, onCancel) => AIDialogRef.value?.show({ type: 'schema', onOk, onCancel }),
    aiPromptDialogRef: AIDialogRef,
  })

  return {
    shareProjectModalRef,
    shareDocModalRef,
    exportDocModalRef,
    AIDialogRef,
  }
}

export function useProjectLayoutContext<K extends keyof ProjectDetailContext>(key: K): Exclude<ProjectDetailContext[K], undefined>
export function useProjectLayoutContext(): ProjectDetailContext
export function useProjectLayoutContext(key?: keyof ProjectDetailContext) {
  const ctx = inject(ProjectDetailContextKey)

  if (key)
    return ctx?.[key]

  return ctx
}
