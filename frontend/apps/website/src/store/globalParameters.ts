import { createGlobalParamerter, deleteGlobalParamerter, getGlobalParamList, updateGlobalParamerter } from '@/api/param'
import { GlobalParameter, GlobalParameters } from '@/typings'
import { defineStore } from 'pinia'

interface GlobalParametersStore {
  parameters: GlobalParameters
}

export const uesGlobalParametersStore = defineStore('definitionParameters', {
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

    async deleteGlobalParameter(project_id: string | number, param: GlobalParameter, is_unref = 1) {
      await deleteGlobalParamerter({ project_id, ...param, is_unref })
    },
  },
})

export const useDefinitionParametersStore = uesGlobalParametersStore

export default uesGlobalParametersStore
