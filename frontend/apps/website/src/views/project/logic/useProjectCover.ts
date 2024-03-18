import { getProjectDefaultCover } from '@/api/project'
import { ProjectListCoverBgColors, ProjectListCoverIcons } from '@/commons'

export function useProjectCover(form: Record<string, any>) {
  const mapper = (value: string) => ({ value, label: value })

  const projectCoverBgColorsOptions = ProjectListCoverBgColors.map(mapper)
  const projectCoverIcons = ProjectListCoverIcons.map(mapper)

  const updateProjectCover = (cover: ProjectAPI.ProjectCover) => {
    form.cover = JSON.stringify(cover)
  }

  const projectCoverRef: ComputedRef<ProjectAPI.ProjectCover> = computed(() => {
    const defaultCover: ProjectAPI.ProjectCover = getProjectDefaultCover()

    if (form.cover) {
      try {
        return JSON.parse(form.cover) as ProjectAPI.ProjectCover
      }
      catch (error) {
        return defaultCover
      }
    }

    // set default cover
    updateProjectCover(defaultCover)
    return defaultCover
  })

  const bgColorRef = computed({
    get() {
      if (projectCoverRef.value.type === 'icon')
        return projectCoverRef.value.coverBgColor

      return null
    },
    set(val) {
      updateProjectCover({ ...projectCoverRef.value, coverBgColor: val as string })
    },
  })

  const iconRef = computed({
    get() {
      if (projectCoverRef.value.type === 'icon')
        return projectCoverRef.value.coverIcon

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
