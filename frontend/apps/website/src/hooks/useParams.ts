import useProjectStore from '@/store/project'

export const useParams = () => {
  const router = useRouter()
  const { projectDetailInfo } = useProjectStore()
  const { iteration_id } = router.currentRoute.value.params

  const params = computed<{ iteration_id: string; doc_id: string; schema_id: string; response_id: string }>(() => {
    const { iteration_id, doc_id, schema_id, response_id } = router.currentRoute.value.params
    return { iteration_id: iteration_id as string, doc_id: doc_id as string, schema_id: schema_id as string, response_id: response_id as string }
  })

  return {
    project_id: projectDetailInfo?.id as string,
    iteration_id: iteration_id as string,
    computedRouteParams: params,
  }
}
