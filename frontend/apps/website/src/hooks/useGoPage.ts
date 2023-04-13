import { getDocumentDetailPath, getSchemaDetailPath, getDocumentEditPath, getSchemaEditPath } from '@/router'

export const useGoPage = () => {
  const router = useRouter()
  const project_id = router.currentRoute.value.params.project_id as string
  const goSchemaDetailPage = (schemaId?: string | number) =>
    router.push(getSchemaDetailPath(project_id as string, schemaId || (router.currentRoute.value.params.shcema_id as string)))
  const goSchemaEditPage = (schemaId?: string | number) => router.push(getSchemaEditPath(project_id as string, schemaId || (router.currentRoute.value.params.shcema_id as string)))

  const goDocumentDetailPage = (docId?: string | number) => router.push(getDocumentDetailPath(project_id as string, docId || (router.currentRoute.value.params.doc_id as string)))
  const goDocumentEditPage = (docId?: string | number) => router.push(getDocumentEditPath(project_id as string, docId || (router.currentRoute.value.params.doc_id as string)))

  return {
    goSchemaDetailPage,
    goSchemaEditPage,

    goDocumentDetailPage,
    goDocumentEditPage,
  }
}
