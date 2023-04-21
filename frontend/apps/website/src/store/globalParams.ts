// import { getProjectList, getProjectDetail, updateProjectBaseInfo, getProjectServerUrlList, saveProjectServerUrlList } from '@/api/param'
import { GlobalParameters } from '@/typings'
import { defineStore } from 'pinia'

interface GlobalParametersStore {
  parameters: GlobalParameters
}

export const uesGlobalParametersStore = defineStore('GlobalParametersStore', {
  state: (): GlobalParametersStore => ({
    parameters: {
      header: [],
      cookie: [],
      query: [],
    },
  }),

  getters: {
    headers: (state) => state.parameters.header,
    cookies: (state) => state.parameters.cookie,
    queries: (state) => state.parameters.query,
  },

  actions: {},
})

export default uesGlobalParametersStore
