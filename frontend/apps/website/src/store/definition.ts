import { getDefinitionList, updateDefinition, createDefinition, copyDefinition } from '@/api/definition'
import { traverseTree } from '@apicat/shared'
import { defineStore } from 'pinia'

export const extendDocTreeFeild = (node = {} as any) => {
  node = node || {}
  node._extend = {
    isLeaf: true,
    isEditable: false,
    isCurrent: false,
  }

  return node
}

export const useDefinitionStore = defineStore('definition', {
  state: () => ({
    definitions: [] as Array<any>,
    tempCreateSchemaParentId: undefined as number | undefined,
  }),
  actions: {
    async getDefinitions(project_id: string) {
      const tree = await getDefinitionList(project_id)
      this.definitions = traverseTree((item: any) => extendDocTreeFeild(item), tree) as any
      return this.definitions
    },

    async updateDefinition(data: any) {
      await updateDefinition(data)
      data.id = data.def_id
      this.updateDefinitionStore(data)
    },

    async createDefinition(data: any) {
      const definition: any = await createDefinition(data)
      this.definitions.unshift(extendDocTreeFeild(definition))
      return definition
    },

    async copyDefinition(project_id: string, def_id: string | number) {
      const definition: any = await copyDefinition(project_id, def_id)
      const index = this.definitions.findIndex((item: any) => item.id === def_id)
      if (index !== -1) {
        this.definitions.splice(index + 1, 0, extendDocTreeFeild(definition))
      }
      return definition
    },

    updateDefinitionStore(definition: any) {
      const { id, name, description, schema } = definition
      const target = this.definitions.find((item: any) => item.id === id)
      target.name = name
      target.description = description
      target.schema = { ...target.schema, ...schema }
    },
  },
})

export default useDefinitionStore
