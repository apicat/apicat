import { defineStore } from 'pinia'
import { ITERATION_DETAIL_PATH_NAME } from '@/router/constant'
import { apiGetIterationInfo } from '@/api/iteration'

interface IterationState {
  iterationInfo: IterationAPI.ResponseIteration | null
}

export const useIterationStore = defineStore('iterationStore', {
  state: (): IterationState => ({
    iterationInfo: null,
  }),

  getters: {
    isIterationRoute(): boolean {
      return !!this.$router.currentRoute.value.matched.find((item) => item.name === ITERATION_DETAIL_PATH_NAME)
    },
  },
  actions: {
    async getIterationInfo(id: string): Promise<IterationAPI.ResponseIteration> {
      this.iterationInfo = await apiGetIterationInfo(id)
      return this.iterationInfo
    },

    gatherIterationInfo(params?: Record<string, any>) {
      params = params || {}
      if (this.isIterationRoute) params.iteration_id = this.$router.currentRoute.value.params.iteration_id

      return params
    },
  },
})
