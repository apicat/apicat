import { storeToRefs } from 'pinia'
import useProjectStore from '@/store/project'

export function onMountedWithProject(callback: (project: ProjectAPI.ResponseProject) => void | Promise<void>): void {
  const { project } = storeToRefs(useProjectStore())
  onMounted(async () => {
    project.value && await callback(project.value)
  })
}
