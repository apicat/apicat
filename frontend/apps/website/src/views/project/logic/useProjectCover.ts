import { ProjectListCoverBgColors, ProjectListCoverIcons, randomArray } from '@/commons'
import { ProjectCover } from '@/typings'

export const getProjectDefaultCover = (overwrite?: Partial<ProjectCover>): ProjectCover => ({
  type: 'icon',
  coverBgColor: randomArray(ProjectListCoverBgColors),
  coverIcon: randomArray(ProjectListCoverIcons),
  ...overwrite,
})

export const useProjectCover = (form: Record<string, any>) => {
  const projectCoverRef: ComputedRef<ProjectCover> = computed(() => {
    const defaultCover: ProjectCover = getProjectDefaultCover()

    if (form.cover) {
      try {
        return JSON.parse(form.cover) as ProjectCover
      } catch (error) {
        return defaultCover
      }
    }

    // set default cover
    updateProjectCover(defaultCover)
    return defaultCover
  })

  const updateProjectCover = (cover: ProjectCover) => {
    form.cover = JSON.stringify(cover)
  }

  const bgColorRef = computed({
    get() {
      if (projectCoverRef.value.type === 'icon') {
        return projectCoverRef.value.coverBgColor
      }

      return null
    },
    set(val) {
      updateProjectCover({ ...projectCoverRef.value, coverBgColor: val as string })
    },
  })

  const iconRef = computed({
    get() {
      if (projectCoverRef.value.type === 'icon') {
        return projectCoverRef.value.coverIcon
      }

      return null
    },
    set(val) {
      updateProjectCover({ ...projectCoverRef.value, coverIcon: val as string })
    },
  })

  return {
    projectCoverRef,
    bgColorRef,
    iconRef,
  }
}
