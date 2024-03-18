import { defineStore } from 'pinia'
import { cloneDeep } from 'lodash-es'
import { addGlobalParameter, deleteGlobalParameter, getGlobalParameters, sortGlobalParameter, updateGlobalParameter } from '@/api/project/setting/globalParameter'

export const useGlobalParameters = defineStore('project.globalParameter', {
  state: () => ({
    parameters: {
      header: [],
      cookie: [],
      query: [],
      path: [],
    } as ProjectAPI.ResponseGlobalParamList,
  }),

  getters: {
    headers: state => state.parameters.header,
    cookies: state => state.parameters.cookie,
    queries: state => state.parameters.query,
  },

  actions: {
    async getGlobalParameterList(projectID: string) {
      this.parameters = await getGlobalParameters(projectID)
      return this.parameters
    },

    async addGlobalParameter(projectID: string, data: Omit<ProjectAPI.GlobalParameter, 'id'>): Promise<ProjectAPI.GlobalParameter> {
      const key = data.in as ProjectAPI.GlobalParameterType
      const parameter = await addGlobalParameter(projectID, cloneDeep(data))
      this.parameters[key].push(parameter)
      return parameter
    },

    async updateGlobalParameter(projectID: string, data: ProjectAPI.GlobalParameter) {
      const key = data.in as ProjectAPI.GlobalParameterType
      await updateGlobalParameter(projectID, cloneDeep(data))
      this.parameters[key] = this.parameters[key].map(p => (p.id === data.id ? data : p))
    },

    async deleteGlobalParameter(projectID: string, { id, ...data }: ProjectAPI.GlobalParameter, deref: boolean = true) {
      const key = data.in as ProjectAPI.GlobalParameterType
      await deleteGlobalParameter(projectID, id, deref)
      this.parameters[key] = this.parameters[key].filter(p => p.id !== id)
    },

    async sortGlobalParameter(projectID: string, data: { oldIndex: number;newIndex: number;in: ProjectAPI.GlobalParameterType }) {
      const key = data.in as ProjectAPI.GlobalParameterType
      const parameters = this.parameters[key]
      if (!parameters || !parameters.length)
        return
      const parameter = parameters.splice(data.oldIndex, 1)[0]
      parameters.splice(data.newIndex, 0, parameter)
      await sortGlobalParameter(projectID, { parameterIDs: parameters.map(p => p.id), in: data.in })
    },

  },
})
