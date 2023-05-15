import { getDefinitionSchemaList, updateDefinitionSchema, createDefinitionSchema, copyDefinitionSchema, deleteDefinitionSchema } from '@/api/definitionSchema'
import { Definition } from '@/components/APIEditor/types'
import { traverseTree } from '@apicat/shared'
import { defineStore } from 'pinia'

export const extendDocTreeFeild = (node = {} as any) => {
  node = node || {}
  node._extend = {
    isLeaf: true,
    isEditable: false,
    isCurrent: false,
  }

  Object.defineProperty(node.schema, '_id', {
    value: node.id,
    enumerable: false,
    writable: false,
    configurable: false,
  })

  return node
}

export const useDefinitionStore = defineStore('definition', {
  state: () => ({
    definitions: [] as Definition[],
    tempCreateSchemaParentId: undefined as number | undefined,
  }),
  actions: {
    async getDefinitions(project_id: string) {
      const tree = await getDefinitionSchemaList(project_id)
      this.definitions = traverseTree((item: any) => extendDocTreeFeild(item), tree) as any
      return this.definitions
    },

    async updateDefinition(data: any) {
      await updateDefinitionSchema(data)
      data.id = data.def_id
      this.updateDefinitionStore(data)
    },

    async createDefinition(data: any) {
      try {
        const definition: any = await createDefinitionSchema(data)
        this.definitions.unshift(extendDocTreeFeild(definition))
        return definition
      } catch (error) {
        // error
      }
    },

    async copyDefinition(project_id: string, def_id: string | number) {
      const definition: any = await copyDefinitionSchema(project_id, def_id)
      const index = this.definitions.findIndex((item: any) => item.id === def_id)
      if (index !== -1) {
        this.definitions.splice(index + 1, 0, extendDocTreeFeild(definition))
      }
      return definition
    },

    updateDefinitionStore(definition: any) {
      const { id, name, description, schema } = definition
      const target = this.definitions.find((item: any) => item.id === id)
      if (target) {
        target.name = name
        target.description = description
        target.schema = { ...target.schema, ...schema }
      }
    },

    async deleteDefinition(project_id: string, def_id: string | number, is_unref = 1) {
      await deleteDefinitionSchema(project_id as string, def_id, is_unref)
    },
  },
})

export default useDefinitionStore
