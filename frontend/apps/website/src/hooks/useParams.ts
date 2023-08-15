import useProjectStore from '@/store/project'

export const useParams = () => {
  const router = useRouter()
  const { projectDetailInfo } = useProjectStore()
  const { doc_id, schema_id, iteration_id, response_id } = router.currentRoute.value.params
  return {
    project_id: projectDetailInfo?.id as string,
    iteration_id: iteration_id as string,
    doc_id: doc_id as string,
    schema_id: schema_id as string,
    response_id: response_id as string,
  }
}
