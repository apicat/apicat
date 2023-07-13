import type { InjectionKey } from 'vue'

export interface ProjectDetailModalsContext {
  exportDocument: (project_id?: string, doc_id?: string | number) => void
  shareDocument: (project_id: string, doc_id: string) => void
  shareProject: (project_id: string) => void
}

export const ProjectDetailModalsContextKey: InjectionKey<ProjectDetailModalsContext> = Symbol('ProjectDetailModalsContextKey')
