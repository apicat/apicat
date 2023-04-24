import { createGlobalParamerter, deleteGlobalParamerter, getGlobalParamList, updateGlobalParamerter } from '@/api/param'
import { GlobalParameter, GlobalParameters } from '@/typings'
import { defineStore } from 'pinia'

interface GlobalParametersStore {
  parameters: GlobalParameters
}

export const uesGlobalParametersStore = defineStore('globalParameters', {
  state: (): GlobalParametersStore => ({
    parameters: {
      header: [],
      cookie: [],
      query: [],
      path: [],
    },
  }),

  getters: {
    headers: (state) => state.parameters.header,
    cookies: (state) => state.parameters.cookie,
    queries: (state) => state.parameters.query,
    paths: (state) => state.parameters.path,
  },

  actions: {
    async getGlobalParameters(project_id: string | number) {
      const data: any = await getGlobalParamList({ project_id })
      this.parameters = data || { header: [], cookie: [], query: [], path: [] }
    },

    async addGlobalParameter(project_id: string | number, type: string, param: GlobalParameter) {
      return await createGlobalParamerter({ project_id, in: type, ...param })
    },

    async updateGlobalParameter(project_id: string | number, type: string, param: GlobalParameter) {
      await updateGlobalParamerter({ project_id, in: type, ...param })
    },

    async deleteGlobalParameter(project_id: string | number, type: string, param: GlobalParameter) {
      await deleteGlobalParamerter({ project_id, in: type, ...param })
    },
  },
})

export default uesGlobalParametersStore
