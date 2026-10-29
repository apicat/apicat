import { createContext, useContext } from 'react'
import type { APIDoc, RequestDraft, Session, Workspace } from './types'
export interface Workbench {
  doc: APIDoc
  workspace: Workspace
  session: Session
  saveWorkspace: (state: Workspace) => Promise<void>
  drafts: Record<string, RequestDraft>
  setDraft: (key: string, draft: RequestDraft) => void
}
export const WorkbenchContext = createContext<Workbench | null>(null)
export function useWorkbench() {
  const ctx = useContext(WorkbenchContext)
  if (!ctx) throw new Error('工作台未加载')
  return ctx
}
