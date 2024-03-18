import type CreateOrUpdateProjectGroup from '@/views/project/components/CreateOrUpdateProjectGroup.vue'

export interface ProjectListContext {
  createOrUpdateProjectGroupRef?: Ref<InstanceType<typeof CreateOrUpdateProjectGroup>>
}

export const ProjectListContextKey: InjectionKey<ProjectListContext> = Symbol('ProjectListContextKey')

export function useProjectListProvider(): ProjectListContext {
  const context = {
    createOrUpdateProjectGroupRef: ref(),
  }

  provide(ProjectListContextKey, context)

  return context
}

export function useProjectListContext(): ProjectListContext | Record<string, never> {
  return inject<ProjectListContext>(ProjectListContextKey) || {}
}
