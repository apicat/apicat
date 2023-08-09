import { ProjectInfo } from '@/typings'
import { Iteration } from '@/typings/iteration'
// import { MemberAuthorityInProject } from '@/typings/member'
import { defineStore } from 'pinia'

interface IterationState {
  // 当前选中所关注的项目Key
  selectedProjectKey: number | null
}

export const useIterationStore = defineStore('iterationStore', {
  state: (): IterationState => ({
    selectedProjectKey: null,
  }),

  actions: {},
})
