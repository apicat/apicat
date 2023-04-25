import { addResponseParam, deleteResponseParam, getCommonResponseList, updateResponseParam } from '@/api/commonResponse'
import { APICatCommonResponse } from '@/typings'

import { defineStore } from 'pinia'

interface CommonResponseState {
  response: APICatCommonResponse[]
}

export const useCommonResponseStore = defineStore('commonResponse', {
  state: (): CommonResponseState => ({
    response: [],
  }),

  getters: {},

  actions: {
    async getCommonResponseList(project_id: string | number) {
      const data: APICatCommonResponse[] = await getCommonResponseList({ project_id })
      this.response = data || []
      return this.response
    },

    async addCommonResponse(project_id: string | number, detail: APICatCommonResponse) {
      const data: any = await addResponseParam({ project_id, ...detail })
      this.response.unshift(data)
      return data
    },

    async updateResponseParam(project_id: string | number, detail: APICatCommonResponse) {
      const { id: response_id, ...rest } = detail
      await updateResponseParam({ project_id, response_id, ...rest })
      const respons = this.response.find((item) => item.id === response_id)
      if (respons) {
        Object.assign(respons, rest)
      }
    },

    async deleteResponseParam(project_id: string | number, detail: APICatCommonResponse) {
      const { id: response_id } = detail
      await deleteResponseParam({ project_id, response_id })
      const index = this.response.findIndex((item) => item.id === response_id)
      index !== -1 && this.response.splice(index, 1)
    },
  },
})

export default useCommonResponseStore
