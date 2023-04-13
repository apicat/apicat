export const useParams = () => {
  const router = useRouter()
  const { project_id, doc_id, shcema_id } = router.currentRoute.value.params
  return { project_id, doc_id, shcema_id }
}
