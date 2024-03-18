import { defineStore } from 'pinia'
import { delay, traverseTree } from '@apicat/shared'
import { apiCopyResponse, apiCreateResponse, apiDeleteResponse, apiEditResponse, apiGetResponseInfo, apiGetResponseTree, apiMoveResponse, apiRenameResponseCategory } from '@/api/project/definition/response'
import { ResponseTypeEnum } from '@/commons'

interface ResponseState {
  loading: boolean
  responseDetail: Definition.ResponseDetail
  responses: Definition.ResponseTreeNode[]
}

function createDefaultResponse(): Definition.ResponseDetail {
  return {
    id: 0,
    name: '',
    parentID: 0,
    type: ResponseTypeEnum.Response,
  }
}

export const useDefinitionResponseStore = defineStore('project.definitionResponse', {
  state: (): ResponseState => ({
    loading: false,
    responseDetail: createDefaultResponse(),
    responses: [],
  }),
  actions: {
    async getResponseDetail(projectID: string, id: number) {
      try {
        // reset response to avoid multiple watch
        this.loading = true
        this.responseDetail = createDefaultResponse()
        this.responseDetail = await apiGetResponseInfo(projectID, id)
      }
      finally {
        this.loading = false
      }
    },
    async getResponses(projectID: string) {
      this.responses = await apiGetResponseTree(projectID)
      return this.responses
    },

    async renameResponseCategory(projectID: string, response: Definition.ResponseTreeNode) {
      return await apiRenameResponseCategory(projectID, response)
    },

    async updateResponse(projectID: string, data: Definition.UpdateResponse) {
      const newData = await apiEditResponse(projectID, data)
      this.updateResponseToStore(data)
      return newData
    },

    async createResponse(projectID: string, data: Definition.CreateResponseTreeNode) {
      return await apiCreateResponse(projectID, data)
    },

    async deleteResponse(projectID: string, id: number, deref: boolean = false) {
      await apiDeleteResponse(projectID, id, deref)
    },

    async copyResponse(projectID: string, id: number) {
      return await apiCopyResponse(projectID, id)
    },

    async moveResponse(projectID: string, data: Definition.RequestSortParams) {
      await apiMoveResponse(projectID, data)
    },

    updateResponseToStore(newResponse: Partial<Definition.UpdateResponse>) {
      traverseTree<Definition.ResponseTreeNode>(
        (response) => {
          if (response.id === newResponse.id) {
            Object.assign(response, newResponse)
            return false
          }
          return true
        },
        this.responses,
        { subKey: 'items' },
      )
    },
  },
})

export default useDefinitionResponseStore
