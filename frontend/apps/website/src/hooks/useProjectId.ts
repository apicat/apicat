export const useProjectId = () => {
  const router = useRouter()
  const project_id = router.currentRoute.value.params.project_id
  return project_id as string | number
}
