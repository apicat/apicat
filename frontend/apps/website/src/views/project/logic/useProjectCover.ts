import { getProjectDefaultCover } from '@/api/project'
import { ProjectListCoverBgColors, ProjectListCoverIcons } from '@/commons'
import { ProjectCover } from '@/typings'

export const useProjectCover = (form: Record<string, any>) => {
  const mapper = (value: string) => ({ value, label: value })

  const projectCoverBgColorsOptions = ProjectListCoverBgColors.map(mapper)
  const projectCoverIcons = ProjectListCoverIcons.map(mapper)

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
    projectCoverBgColorsOptions,
    projectCoverIcons,

    projectCoverRef,
    bgColorRef,
    iconRef,
  }
}
