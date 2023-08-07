import { getDocumentDetailPath, getSchemaDetailPath, getDocumentEditPath, getSchemaEditPath, getDefinitionResponseDetailPath, getDefinitionResponseEditPath } from '@/router'

export const useGoPage = () => {
  const router = useRouter()
  const project_id = router.currentRoute.value.params.project_id as string

  const goSchemaDetailPage = (schemaId?: string | number) =>
    router.push(getSchemaDetailPath(project_id as string, schemaId || (router.currentRoute.value.params.schema_id as string)))
  const goSchemaEditPage = (schemaId?: string | number) => router.push(getSchemaEditPath(project_id as string, schemaId || (router.currentRoute.value.params.schema_id as string)))

  const goResponseDetailPage = (responseId?: string | number) =>
    router.push(getDefinitionResponseDetailPath(project_id as string, responseId || (router.currentRoute.value.params.response_id as string)))
  const goResponseEditPage = (responseId?: string | number) =>
    router.push(getDefinitionResponseEditPath(project_id as string, responseId || (router.currentRoute.value.params.response_id as string)))

  const goDocumentDetailPage = (docId?: string | number) => router.push(getDocumentDetailPath(project_id as string, docId || (router.currentRoute.value.params.doc_id as string)))
  const goDocumentEditPage = (docId?: string | number) => router.push(getDocumentEditPath(project_id as string, docId || (router.currentRoute.value.params.doc_id as string)))

  return {
    goSchemaDetailPage,
    goSchemaEditPage,

    goDocumentDetailPage,
    goDocumentEditPage,

    goResponseDetailPage,
    goResponseEditPage,
  }
}
