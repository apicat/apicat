import { getIterationDetail } from '@/api/iteration'
import { ITERATION_DETAIL_PATH_NAME } from '@/router/constant'
import { Iteration } from '@/typings'
import { defineStore } from 'pinia'

interface IterationState {
  iterationInfo: Iteration | null
}

export const useIterationStore = defineStore('iterationStore', {
  state: (): IterationState => ({
    iterationInfo: null,
  }),

  getters: {
    isIterationRoute: () => {
      const router = useRouter()
      return !!router.currentRoute.value.matched.find((item) => item.name === ITERATION_DETAIL_PATH_NAME)
    },
  },
  actions: {
    async getIterationInfo(iteration_id: string): Promise<Iteration> {
      const iterationDetail = await getIterationDetail({ iteration_id })
      this.iterationInfo = iterationDetail
      return iterationDetail
    },
  },
})
