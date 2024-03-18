import { storeToRefs } from 'pinia'
import useProjectStore from '@/store/project'
import router from '@/router'

export function useParams() {
  const { projectID } = storeToRefs(useProjectStore())
  return {
    project_id: projectID.value as string,
    projectID: projectID.value as string,
    iterationID: computed(() => router.currentRoute.value.params.iterationID as string),
  }
}
