export const useParams = () => {
  const router = useRouter()
  const { project_id, doc_id, shcema_id } = router.currentRoute.value.params
  return { project_id: project_id as string, doc_id: doc_id as string, shcema_id: shcema_id as string }
}
