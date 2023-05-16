import { getDefinitionResponseList, updateDefinitionResponse, createDefinitionResponse, deleteDefinitionResponse } from '@/api/definitionResponse'
import { DefinitionTypeEnum } from '@/commons'
import { DefinitionResponse } from '@/typings'
import { traverseTree } from '@apicat/shared'
import { defineStore } from 'pinia'

export const extendDocTreeFeild = (node = {} as any) => {
  node = node || {}
  node._extend = {
    isLeaf: node.type !== DefinitionTypeEnum.DIR,
    isEditable: false,
    isCurrent: false,
  }
  return node
}

export const useDefinitionResponseStore = defineStore('definitionResponse', {
  state: () => ({
    responses: [] as DefinitionResponse[],
  }),
  actions: {
    async getDefinitions(project_id: string) {
      const tree = await getDefinitionResponseList(project_id)
      this.responses = traverseTree((item: any) => extendDocTreeFeild(item), tree) as any
      return this.responses
    },

    async updateDefinition(data: DefinitionResponse) {
      await updateDefinitionResponse(data)
      this.updateDefinitionStore(data)
    },

    async createDefinition(data: any) {
      try {
        const definition: any = await createDefinitionResponse(data)
        this.responses.unshift(extendDocTreeFeild(definition))
        return definition
      } catch (error) {
        // error
      }
    },

    updateDefinitionStore(definition: DefinitionResponse) {
      const { id, name, description, content, header } = definition
      const target = this.responses.find((item: DefinitionResponse) => item.id === id)
      if (target) {
        target.name = name
        target.description = description
        target.header = header
        target.content = content
      }
    },

    async deleteDefinition(project_id: string, id: string | number, is_unref = 1) {
      await deleteDefinitionResponse(project_id as string, id, is_unref)
    },
  },
})

export default useDefinitionResponseStore
